# smartbank CSV exporter

## status
Project is complete and works for current version of www.banksmart.pl
However it was tested only with standard and savings accounts in PLN.

## polska wersja
Ten program służy do eksportu danych z www.banksmart.pl do formatu CSV.
Powstał w celu, aby nie wykonywać tego ręcznie dla każdego okresu.

Ponieważ przetwarza istotne dane dostępu do konta,
jest udostępniony w postaci źródeł.

Instrukcja instalacji `go` jest tutaj: https://golang.org/dl/
następnie wystarczy wykonać:

`go get -u github.com/Komosa/smartbank`
a następnie można uruchomić program.

Jako wynik powstanie plik `output.csv`
