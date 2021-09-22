package wallet

import (
	"bufio"
	"errors"
	"github.com/google/uuid"
	"github.com/rustamfozilov/wallet/pkg/types"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

//var ErrAccountNotFound = errors.New("account not found")
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil

}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	var account *types.Account

	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return ErrAccountNotFound
	}
	account.Balance += amount
	return nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

//func (s *Service) Reject(paymentID string) error {
//	for _, payment := range s.payments {
//		if payment.ID == paymentID {
//			payment.Status = types.PaymentStatusFail
//			_, err := s.FindAccountByID(payment.AccountID)
//			if err != nil {
//				return err
//			}
//			//acc.Balance = acc.Balance - payment.Amount
//			return nil
//		}
//	}
//	return ErrPaymentNotFound
//}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	//log.Print("[")
	//for _, payment := range s.payments {
	//	log.Print(*payment)
	//}
	//log.Println("]")
	//log.Println("paymentID:", paymentID)
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	payment := &types.Payment{
		ID:        uuid.New().String(),
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	// to do acc
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, err
	}
	account.Balance = account.Balance - payment.Amount
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}
	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, err
	}
	var repeatedPayment = types.Payment{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Category:  payment.Category,
		Status:    payment.Status,
	}
	//log.Println("reapetedPayment",repeatedPayment)
	s.payments = append(s.payments, &repeatedPayment)
	account.Balance = account.Balance - payment.Amount
	return &repeatedPayment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	//log.Print("[")
	//for _, payment := range s.payments {
	//	log.Print(*payment)
	//}
	//log.Println("]")
	//log.Println("paymentID:", paymentID)
	var favorite = types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, &favorite)
	return &favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {

	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}
	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}

func (s *
Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {

	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	for _, account := range s.accounts {
		line := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" + strconv.FormatInt(int64(account.Balance), 10) + "|"
		_, err = file.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	accounts, err2 := readAll(file)
	if err2 != nil {
		return err2
	}
	lines := strings.Split(string(accounts), "|")
	if len(lines) == 0 {
		return err
	}
	lines = lines[:len(lines)-1]
	log.Println(lines)
	for _, line := range lines {

		fields := strings.Split(line, ";")
		log.Println(fields)
		idString := fields[0]
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			return err
		}
		phone := types.Phone(fields[1])
		balanceString := fields[2]
		balance, err := strconv.ParseInt(balanceString, 10, 64)
		if err != nil {
			return err
		}
		var acc = types.Account{
			ID:      id,
			Phone:   phone,
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, &acc)
	}
	return nil
}

func readAll(reader io.Reader) ([]byte, error) {
	accounts := make([]byte, 0)
	buf := make([]byte, 4096)
	for {
		read, err := reader.Read(buf)
		if err == io.EOF {
			accounts = append(accounts, buf[:read]...)
			break
		}
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, buf[:read]...)
	}
	return accounts, nil
}

func (s *Service) Export(dir string) error {
	err := exportAccounts(dir, s)
	if err != nil {
		return err
	}
	err = exportPayments(dir, s)
	if err != nil {
		return err
	}
	return exportFavorites(dir, s)
}

func exportFavorites(dir string, s *Service) error {
	if len(s.favorites) == 0 {
		return nil
	}
	var line string
	for _, favorite := range s.favorites {
		line += favorite.ID + "|" + strconv.FormatInt(int64(favorite.AccountID), 10) + "|" +
			favorite.Name + "|" + strconv.FormatInt(int64(favorite.Amount), 10) +
			"|" + string(favorite.Category) + "\n"
	}
	err := ioutil.WriteFile(path.Join(dir, "favorites.dump"), []byte(line), 0666)
	if err != nil {
		return err
	}
	return nil
}

func exportPayments(dir string, s *Service) error {
	if len(s.payments) == 0 {
		return nil
	}
	var line string
	for _, payment := range s.payments {
		line = creatingLine(line, payment)
	}
	err := ioutil.WriteFile(path.Join(dir, "payments.dump"), []byte(line), 0666)
	if err != nil {
		return err
	}
	return nil
}

func creatingLine(line string, payment *types.Payment) string {
	line += payment.ID + "|" + strconv.FormatInt(payment.AccountID, 10) + "|" +
		strconv.FormatInt(int64(payment.Amount), 10) + "|" + string(payment.Category) + "|" +
		string(payment.Status) + "\n"
	return line
}

func exportAccounts(dir string, s *Service) error {
	if len(s.accounts) == 0 {
		return nil
	}
	var line string
	for _, account := range s.accounts {
		line += strconv.FormatInt(account.ID, 10) + "|" + string(account.Phone) +
			"|" + strconv.FormatInt(int64(account.Balance), 10) + "\n"
	}

	err := ioutil.WriteFile(path.Join(dir, "accounts.dump"), []byte(line), 0666)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Import(dir string) error {
	err := s.ImportAccounts(dir)
	if err != nil {
		return err
	}
	err = s.ImportPayments(dir)
	if err != nil {
		return err
	}
	err = s.ImportFavorites(dir)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ImportAccounts(dir string) error {
	fileAccounts, err := os.Open(path.Join(dir, "accounts.dump"))
	if err != nil {
		return nil
	}
	defer func() {
		if err2 := fileAccounts.Close(); err2 != nil {
			log.Println(err2)
		}
	}()
	reader := bufio.NewReader(fileAccounts)
	var lines = make([]string, 0)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Println("line:", line)
		line = line[:len(line)-len("\n")]
		lines = append(lines, line)
		log.Println("lines:", lines)
	}
	for _, line := range lines {
		accFromFile, err := parseAccountLine(line)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("accfromfile:", accFromFile)
		if index := s.accountExists(accFromFile.ID); index != 0 { // update
			s.accounts[index] = accFromFile
			continue
		}
		s.accounts = append(s.accounts, accFromFile) // add
		s.nextAccountID++
	}
	return nil
}

func parseAccountLine(line string) (*types.Account, error) {
	fields := strings.Split(line, "|")
	if len(fields) < 3 {
		return nil, errors.New("wrong line format")
	}
	id, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return nil, err
	}
	balance, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return nil, err
	}
	return &types.Account{
		ID:      id,
		Phone:   types.Phone(fields[1]),
		Balance: types.Money(balance),
	}, nil
}

func (s *Service) accountExists(id int64) int {
	for index, account := range s.accounts {
		if account.ID == id {
			return index
		}
	}
	return 0
}

func (s *Service) ImportPayments(dir string) error {
	filePayments, err := os.Open(path.Join(dir, "payments.dump"))
	if err != nil {
		return nil
	}
	defer func() {
		if err2 := filePayments.Close(); err2 != nil {
			log.Println(err2)
		}
	}()
	reader := bufio.NewReader(filePayments)
	var lines = make([]string, 0)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Println("line:", line)
		line = line[:len(line)-len("\n")]
		lines = append(lines, line)
		log.Println("lines:", lines)
	}
	for _, line := range lines {
		paymentFromFile, err := parsePaymentLine(line)
		if err != nil {
			log.Println(err)
			continue
		}
		if index := s.paymentExists(paymentFromFile.ID); index != 0 {
			s.payments[index] = paymentFromFile
			continue
		}
		s.payments = append(s.payments, paymentFromFile)
	}
	return nil
}

func parsePaymentLine(line string) (*types.Payment, error) {
	fields := strings.Split(line, "|")
	log.Println("fields:", fields)

	accountID, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return nil, err
	}
	amount, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return nil, err
	}
	return &types.Payment{
		ID:        fields[0],
		AccountID: accountID,
		Amount:    types.Money(amount),
		Category:  types.PaymentCategory(fields[3]),
		Status:    types.PaymentStatus(fields[4]),
	}, nil
}
func (s *Service) paymentExists(id string) int {
	for index, payment := range s.payments {
		if payment.ID == id {
			return index
		}
	}
	return 0
}

func (s *Service) ImportFavorites(dir string) error {
	fileFavorites, err := os.Open(path.Join(dir, "favorites.dump"))
	if err != nil {
		return nil
	}
	defer func() {
		if err2 := fileFavorites.Close(); err2 != nil {
			log.Println(err2)
		}
	}()
	reader := bufio.NewReader(fileFavorites)
	var lines = make([]string, 0)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Println("line:", line)
		line = line[:len(line)-len("\n")]
		lines = append(lines, line)
		log.Println("lines:", lines)
	}
	for _, line := range lines {
		favoriteFromFile, err := parseFavoriteLine(line)
		if err != nil {
			log.Println(err)
			continue
		}
		if index := s.favoriteExists(favoriteFromFile.ID); index != 0 {
			s.favorites[index] = favoriteFromFile
			continue
		}
		s.favorites = append(s.favorites, favoriteFromFile)
	}
	return nil

}

func parseFavoriteLine(line string) (*types.Favorite, error) {
	fields := strings.Split(line, "|")
	log.Println("fields fav:", fields)
	accountID, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return nil, err
	}
	amount, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return nil, err
	}
	return &types.Favorite{
		ID:        fields[0],
		AccountID: accountID,
		Name:      fields[2],
		Amount:    types.Money(amount),
		Category:  types.PaymentCategory(fields[4]),
	}, nil
}

func (s *Service) favoriteExists(id string) int {
	for index, favorite := range s.favorites {
		if favorite.ID == id {
			return index
		}
	}
	return 0
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	accountsPayments := make([]types.Payment, 0)
	for _, payment := range s.payments {
		if payment.AccountID == account.ID {
			accountsPayments = append(accountsPayments, *payment)
		}
	}
	return accountsPayments, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if len(payments) == 0 {
		return nil
	}
	log.Println("payments:", payments, "dir:", dir)
	log.Println("len paymens:", len(payments), "records:", records)
	if len(payments) <= records {
		var line string
		for _, payment := range payments {
			line = creatingLine(line, &payment)
		}
		log.Println("line", line)
		filename := "payments.dump"
		err := ioutil.WriteFile(path.Join(dir, filename), []byte(line), 0666)
		if err != nil {
			return err
		}
		return nil
	}

	fileN := len(payments) / records
	if len(payments)%records != 0 {
		fileN++
	}

	for numberFile := 1; numberFile <= fileN-1; numberFile++ {
		var line string
		for i := 0; i < records; i++ {
			line = creatingLine(line, &payments[i])
		}
		payments = payments[records:]
		log.Println("line2:", line)
		filename := "payments" + strconv.Itoa(numberFile) + ".dump"
		log.Println("filename:", filename)
		err := ioutil.WriteFile(path.Join(dir, filename), []byte(line), 0666)
		if err != nil {
			return err
		}
	}

	line := ""
	for _, payment := range payments {
		line = creatingLine(line, &payment)
	}
	log.Println("line2:", line)
	filename := "payments" + strconv.Itoa(fileN) + ".dump"
	log.Println("filename:", filename)
	err := ioutil.WriteFile(path.Join(dir, filename), []byte(line), 0666)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Hf(payments []types.Payment, dir string, records int) error {

	for numberFile := 1; numberFile <= (len(payments)/records)+1; numberFile++ {
		var line string
		for i := 0; i < records; i++ {
			line = creatingLine(line, &payments[i])
		}
		payments = payments[records:]
		log.Println("line2:", line)
		filename := "payments" + strconv.Itoa(numberFile) + ".dump"
		log.Println("filename:", filename)
		err := ioutil.WriteFile(path.Join(dir, filename), []byte(line), 0666)
		if err != nil {
			return err
		}
	}

	var index int
	numberFile := 1
	for numberFile = 1; numberFile <= len(payments)/records; numberFile++ {
		var line string
		for i := 0; i < records; i++ {
			payment := payments[index]
			line = creatingLine(line, &payment)
			index++

		}
		log.Println("line2:", line)
		filename := "payments" + strconv.Itoa(numberFile) + ".dump"
		log.Println("filename:", filename)
		err := ioutil.WriteFile(path.Join(dir, filename), []byte(line), 0666)
		if err != nil {
			return err
		}
		numberFile++
	}
	var line string
	for i := 0; i < records; i++ {
		payment := payments[index]
		line = creatingLine(line, &payment)
		index++
	}
	filename := "payments" + strconv.Itoa(numberFile) + ".dump"
	err := ioutil.WriteFile(path.Join(dir, filename), []byte(line), 0666)
	if err != nil {
		return err
	}
	return nil

}
