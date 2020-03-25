package main

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (a Account) defaultPassword() string {
	return fmt.Sprintf(`%s%02d`, a.Class.Name, a.Number)
}

func generatePassword(password string) []byte {
	if password == "" {
		return []byte{}
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return pwd
}

func (a Account) login(password string) bool {
	fmt.Println(password)
	fmt.Println(a.Password)
	fmt.Println(a.defaultPassword())
	if a.EmptyPassword() {
		return a.defaultPassword() == password
	}
	return bcrypt.CompareHashAndPassword(a.Password, []byte(password)) == nil
}

func (a Account) EmptyPassword() bool {
	return len(a.Password) == 0
}

// 重設密碼
func (h Handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	acct, _ := session(r)
	account, err := h.getAccount(r.FormValue("account_id"))
	if err == AccountNotFound {
		account = acct
	}

	if !accountPermission(acct, account) {
		w.WriteHeader(403)
		h.resetPage(w, addMessage(r, "您沒有權限變更 "+account.Name+" 的密碼"))
		return
	}

	current := r.FormValue("current_password")
	new := r.FormValue("new_password")
	fmt.Println(current, new)
	if acct.ID == account.ID {
		err = account.resetSelfPassword(current, new)
	} else {
		if acct.login(current) {
			account.Password = []byte{}
		} else {
			err = WrongPassword
		}
	}

	fmt.Println(account.Password)

	if err == WrongPassword {
		h.resetPage(w, addMessage(r, WrongPassword.Error()))
		return
	}

	if err := h.db.Save(&account).Error; err != nil {
		w.WriteHeader(500)
		h.resetPage(w, addMessage(r, err.Error()))
		return
	}

	h.resetPage(w, addMessage(r, "已重設 "+account.Name+" 的密碼"))
}

func (a *Account) resetSelfPassword(current, new string) error {
	if !a.EmptyPassword() && !a.login(current) {
		return WrongPassword
	}
	a.Password = generatePassword(new)
	return nil
}
