package seed

import (
	"fmt"
	"gilsaputro/dating-apps/internal/store/user"
	"gilsaputro/dating-apps/models"
	"gilsaputro/dating-apps/pkg/hash"
)

const numUser = 20

func GenerateSeed(store user.UserStoreMethod, hash hash.HashMethod) error {
	count, err := store.Count()
	if err != nil {
		return err
	}

	newHash, _ := hash.HashValue("banana1")

	if count <= 0 {
		for i := 0; i <= numUser; i++ {
			// Create a new faker instance
			store.CreateUser(models.User{
				Username:   fmt.Sprintf("username_%v", i),
				Fullname:   fmt.Sprintf("User Person %v", i),
				Password:   string(newHash),
				Email:      fmt.Sprintf("email%v@fake.com", i),
				IsVerified: false,
			})
		}
	}
	return nil
}
