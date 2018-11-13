package itest

import (
	"fmt"
	"math/rand"

	"github.com/iost-official/go-iost/ilog"
)

// ITest is the test controller
type ITest struct {
	bank    *Account
	keys    []*Key
	clients []*Client
}

// New will return the itest by config and keys
func New(c *Config, keys []*Key) *ITest {
	return &ITest{
		bank:    c.Bank,
		keys:    keys,
		clients: c.Clients,
	}
}

// Load will load the itest from file
func Load(keysfile, configfile string) (*ITest, error) {
	ilog.Infof("Load itest from file...")

	keys, err := LoadKeys(keysfile)
	if err != nil {
		return nil, fmt.Errorf("load keys failed: %v", err)
	}

	itc, err := LoadConfig(configfile)
	if err != nil {
		return nil, fmt.Errorf("load itest config failed: %v", err)
	}

	it := New(itc, keys)

	ilog.Infof("Load itest from file successful!")
	return it, nil
}

// CreateAccountN will create n accounts concurrently
func (t *ITest) CreateAccountN(num int) ([]*Account, error) {
	ilog.Infof("Create %v account...", num)

	var res chan interface{}
	for i := 0; i < num; i++ {
		go func(n int, res chan interface{}) {
			name := fmt.Sprintf("account%04d", n)
			account, err := t.CreateAccount(name)
			if err != nil {
				res <- err
			} else {
				res <- account
			}
		}(i, res)
	}

	accounts := []*Account{}
	for i := 0; i < num; i++ {
		select {
		case r := <-res:
			err, ok := r.(error)
			if ok {
				ilog.Errorf("Create account failed: %v", err)
				break
			}
			account, ok := r.(*Account)
			if ok {
				accounts = append(accounts, account)
				break
			}
		}
	}

	ilog.Infof("Create %v account successful!", len(accounts))

	if len(accounts) != num {
		return nil, fmt.Errorf("create %v account failed", num-len(accounts))
	}

	// TODO Get account by rpc, and compare account result

	return accounts, nil
}

// CreateAccount will create a account by name
func (t *ITest) CreateAccount(name string) (*Account, error) {
	if len(t.keys) == 0 {
		return nil, fmt.Errorf("keys is empty")
	}
	if len(t.clients) == 0 {
		return nil, fmt.Errorf("clients is empty")
	}
	kIndex := rand.Intn(len(t.keys)) // nolint: golint
	key := t.keys[kIndex]
	cIndex := rand.Intn(len(t.clients))
	client := t.clients[cIndex]

	account, err := client.CreateAccount(t.bank, name, key)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// Transfer will transfer token from sender to recipient
func (t *ITest) Transfer(sender *Account, token, recipient, amount string) error {
	cIndex := rand.Intn(len(t.clients))
	client := t.clients[cIndex]

	err := client.Transfer(sender, token, recipient, amount)
	if err != nil {
		return err
	}

	return nil
}

// SetContract will set the contract on blockchain
func (t *ITest) SetContract(contract *Contract) error {
	cIndex := rand.Intn(len(t.clients))
	client := t.clients[cIndex]

	err := client.SetContract(t.bank, contract)
	if err != nil {
		return err
	}

	return nil
}

// GetTransaction will get transaction by tx hash
func (t *ITest) GetTransaction(hash string) (*Transaction, error) {
	cIndex := rand.Intn(len(t.clients))
	client := t.clients[cIndex]

	transaction, err := client.GetTransaction(hash)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetAccount will get account by name
func (t *ITest) GetAccount(name string) (*Account, error) {
	cIndex := rand.Intn(len(t.clients))
	client := t.clients[cIndex]

	account, err := client.GetAccount(name)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// SendTransaction will send transaction to blockchain
func (t *ITest) SendTransaction(transaction *Transaction) (string, error) {
	cIndex := rand.Intn(len(t.clients))
	client := t.clients[cIndex]

	hash, err := client.SendTransaction(transaction)
	if err != nil {
		return "", err
	}

	return hash, nil
}