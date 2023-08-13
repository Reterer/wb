package utils

import (
	"math/rand"
	"strconv"
	"time"
)

type Bank interface {
	Add(uid string, amount int)
	Check(uid string) int
}

func genUids(n int) []string {
	uids := make([]string, 0, n)

	for i := 0; i < n; i++ {
		uids = append(uids, "test_"+strconv.Itoa(i))
	}
	return uids
}

// Настройки для StressBank ()
type StressBankArgs struct {
	Timeout      time.Duration // timeout Количество времени, которое будет занимать тест
	AddWorkers   int           // addWorkers Количество горутин, которые будут изменять значения
	CheckWorkers int           // checkWorkers Количетсво горутин, которые будут проверять значения
	Users        int           // users Количество уникальных пользователей
}

// Генерирует нагрузку на bank

func StressBank(bank Bank, args StressBankArgs) {
	timeout := args.Timeout
	addWorkers := args.AddWorkers
	checkWorkers := args.CheckWorkers
	users := args.Users

	// Создадим пользователей
	userUids := genUids(users)

	// Контроль горутин
	done := make(chan struct{})

	// Запускаем addWorkers
	for i := 0; i < addWorkers; i++ {
		go func() {
			for {
				// Делаем неблокирующий select
				select {
				case <-done:
					return
				default:
				}
				// Условно берем пользователя какого-нибудь
				user := userUids[rand.Intn(len(userUids))]
				bank.Add(user, rand.Intn(50)-25)
			}
		}()
	}
	// Запускаем checkWorkers
	for i := 0; i < checkWorkers; i++ {
		go func() {
			for {
				// Делаем неблокирующий select
				select {
				case <-done:
					return
				default:
				}
				// Условно берем пользователя какого-нибудь
				user := userUids[rand.Intn(len(userUids))]
				bank.Check(user)
			}
		}()
	}

	time.Sleep(timeout)
	close(done)
}
