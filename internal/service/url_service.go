package service

import (
	"crypto/sha256"
	"errors"
	"url-shortener/pkg/storage"
)

const Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
const hashLength = 8                   // Длина желаемого хэша
const maxIndex = len(Alphabet) * len(Alphabet) - 1 // Максимальное кол-во коллизий для одного хэша

type Storage interface {
	GetLongUrl(shortUrl string) (string, error)
	Insert(shortUrl, longUrl string) error
}

type ShortenerService struct {
	Storage Storage
}

func NewShortenerService(Storage Storage) *ShortenerService {
	return &ShortenerService{Storage: Storage}
}

func (s ShortenerService) Shortening(longUrl string) (string, error) {
	id := 0
	hash := EncodeHash(longUrl)
	for id < maxIndex {
		shortUrl := hash + IntToIndex63(id)
		longCheck, err := s.Storage.GetLongUrl(shortUrl)
		if err == storage.ErrNotFound {
			err = s.Storage.Insert(shortUrl, longUrl)
			return shortUrl, err
		} else if longCheck == longUrl {
			return shortUrl, err
		}
		id++
	}
	return "", errors.New("index out of range")
}

func (s ShortenerService) Expansion(shortUrl string) (string, error) {
	res, err := s.Storage.GetLongUrl(shortUrl)
	return res, err
}

// Для разрешения коллизий вычисляем хэш полученного значения,
// после чего преобразуем его в 8 символов из требуемого алфавита (буквы, цифры, "_")
// в оставшиеся два символа записываем порядковый номер данной коллизий в 63-ом формате

// т.е: предположим что для ссылки "example1.com" вычислился хэш - "Xa31kJi_", и для него нет коллизий - в результате получится значение "Xa31kJi_00"
// предположим, что что для ссылки "example2.com" получился такой же хэш (вероятность данного события крайне мала), в базе уже присутсвует ссылка
// начинающаяся на "Xa31kJi_", тогда результатом для "example2.com" будет - "Xa31kJi_01"

// Функция для преобразования байтов в строку фиксированной длины
func EncodeHash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)

	result := ""
	for i := 0; i < hashLength; i++ {
		index := int(hash[i]) % len(Alphabet) // Получаем индекс из хэша
		result += string(Alphabet[index])     // Добавляем символ из алфавита
	}
	return result
}

// Функция для преобразования числа из 10-тичной системы счисления в двухразрядное число в 63-ной системе
func IntToIndex63(id int) string {
	res := ""
	res += string(Alphabet[id/len(Alphabet)])
	res += string(Alphabet[id%len(Alphabet)])
	return res

}
