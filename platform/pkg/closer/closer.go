package closer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Closer struct {
	mu     sync.Mutex                        // Синхронизация при добавлении функций закрытия
	once   sync.Once                         // Выполнение одного действия разово
	funcs  []func(ctx context.Context) error // Функции закрытия
	logger Logger
}

func New(l Logger) *Closer {
	return &Closer{
		logger: l,
	}
}

func (c *Closer) SetLogger(l Logger) {
	c.logger = l
}

func (c *Closer) AddNamed(name string, f func(ctx context.Context) error) {
	c.add(
		func(ctx context.Context) error {
			start := time.Now()
			c.logger.Info(ctx, fmt.Sprintf("🧩 Закрываем %s...", name))

			err := f(ctx)

			duration := time.Since(start)
			if err != nil {
				c.logger.Error(ctx, fmt.Sprintf("❌ Ошибка при закрытии %s: %v (заняло %s)", name, err, duration))
			} else {
				c.logger.Info(ctx, fmt.Sprintf("✅ %s успешно закрыт за %s", name, duration))
			}

			return err
		},
	)
}

func (c *Closer) CloseAll(ctx context.Context) error {
	var shutdownErrors []error

	c.once.Do(
		func() {
			// Очистка памяти из структуры
			c.mu.Lock()
			funcs := c.funcs
			c.funcs = nil
			c.mu.Unlock()

			if len(funcs) == 0 {
				c.logger.Info(ctx, "ℹ️ Нет функций для закрытия.")
				return
			}

			errCh := make(chan error, len(funcs))
			var wg sync.WaitGroup

			// Выполняем в обратном порядке добавления
			for i := len(funcs) - 1; i >= 0; i-- {
				f := funcs[i]
				wg.Add(1)
				go func(f func(ctx context.Context) error) {
					defer wg.Done()

					// Защита от паники
					defer func() {
						if r := recover(); r != nil {
							errCh <- fmt.Errorf("panic recovered in closer")
							c.logger.Error(ctx, "⚠️ Panic в функции закрытия", zap.Any("error", r))
						}
					}()

					if err := f(ctx); err != nil {
						errCh <- err
					}
				}(f)
			}

			// Закрываем канал ошибок, когда все функции завершатся
			go func() {
				wg.Wait()
				close(errCh)
			}()

			// Обработка ошибок при закрытии
			for {
				select {
				case <-ctx.Done():
					c.logger.Info(ctx, "⚠️ Контекст отменён во время закрытия", zap.Error(ctx.Err()))
					shutdownErrors = append(shutdownErrors, ctx.Err())
					return
				case err, ok := <-errCh:
					if !ok {
						c.logger.Info(ctx, "✅ Все ресурсы успешно закрыты")
						return
					}
					shutdownErrors = append(shutdownErrors, err)
				}
			}
		},
	)

	return errors.Join(shutdownErrors...)
}

func (c *Closer) add(f ...func(ctx context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}
