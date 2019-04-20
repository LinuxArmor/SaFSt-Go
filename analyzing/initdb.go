package analyzing

import (
	"github.com/syndtr/goleveldb/leveldb"
)

const PERMFOLDER = "/usr/local/var/safst/perm"

var perm = map[string]string{
	"abort@libc.so.6":     "0",
	"abs@libc.so.6":       "0",
	"atof@libc.so.6":      "0",
	"atoi@libc.so.6":      "0",
	"atol@libc.so.6":      "0",
	"calloc@libc.so.6":    "8",
	"clearerr@libc.so.6":  "0",
	"delay@libc.so.6":     "0",
	"div@libc.so.6":       "0",
	"exit@libc.so.6":      "0",
	"fclose@libc.so.6":    "8",
	"feof@libc.so.6":      "0",
	"ferror@libc.so.6":    "0",
	"fets@libc.so.6":      "4",
	"fflush@libc.so.6":    "8",
	"fgetc@libc.so.6":     "4",
	"fgetpos@libc.so.6":   "0",
	"fgets@libc.so.6":     "4",
	"fgetwc@libc.so.6":    "4",
	"fgetws@libc.so.6":    "4",
	"fopen@libc.so.6":     "12",
	"fprintf@libc.so.6":   "8",
	"fputc@libc.so.6":     "8",
	"fputs@libc.so.6":     "8",
	"fputwc@libc.so.6":    "8",
	"fputws@libc.so.6":    "8",
	"fread@libc.so.6":     "4",
	"free@libc.so.6":      "32",
	"freopen@libc.so.6":   "12",
	"fscanf@libc.so.6":    "4",
	"fseek@libc.so.6":     "0",
	"fseeko@libc.so.6":    "0",
	"fsetpos@libc.so.6":   "0",
	"ftell@libc.so.6":     "0",
	"ftello@libc.so.6":    "0",
	"fwide@libc.so.6":     "0",
	"fwprintf@libc.so.6":  "8",
	"fwrite@libc.so.6":    "8",
	"fwscanf@libc.so.6":   "4",
	"getc@libc.so.6":      "4",
	"getchar@libc.so.6":   "4",
	"getenv@libc.so.6":    "4",
	"gets@libc.so.6":      "4",
	"getwc@libc.so.6":     "4",
	"getwchar@libc.so.6":  "4",
	"malloc@libc.so.6":    "8",
	"perror@libc.so.6":    "0",
	"printf@libc.so.6":    "8",
	"putc@libc.so.6":      "8",
	"putchar@libc.so.6":   "8",
	"putenv@libc.so.6":    "8",
	"puts@libc.so.6":      "8",
	"putwc@libc.so.6":     "8",
	"putwchar@libc.so.6":  "8",
	"rand@libc.so.6":      "0",
	"realloc@libc.so.6":   "8",
	"remove@libc.so.6":    "32",
	"rename@libc.so.6":    "8",
	"rewind@libc.so.6":    "0",
	"scanf@libc.so.6":     "4",
	"setbuf@libc.so.6":    "0",
	"setenv@libc.so.6":    "8",
	"setvbuf@libc.so.6":   "0",
	"snprintf@libc.so.6":  "8",
	"sprintf@libc.so.6":   "8",
	"sscanf@libc.so.6":    "4",
	"strod@libc.so.6":     "0",
	"strtod@libc.so.6":    "0",
	"strtol@libc.so.6":    "0",
	"swprintf@libc.so.6":  "8",
	"swscanf@libc.so.6":   "4",
	"system@libc.so.6":    "63",
	"tmpfile@libc.so.6":   "12",
	"tmpnam@libc.so.6":    "0",
	"ungetc@libc.so.6":    "8",
	"ungetwc@libc.so.6":   "8",
	"vfprintf@libc.so.6":  "8",
	"vfscanf@libc.so.6":   "4",
	"vfwprintf@libc.so.6": "8",
	"vfwscanf@libc.so.6":  "4",
	"vprintf@libc.so.6":   "8",
	"vscanf@libc.so.6":    "4",
	"vsnprintf@libc.so.6": "8",
	"vsprintf@libc.so.6":  "8",
	"vsscanf@libc.so.6":   "4",
	"vswprintf@libc.so.6": "8",
	"vswscanf@libc.so.6":  "4",
	"vwprintf@libc.so.6":  "8",
	"vwscanf@libc.so.6":   "4",
	"wprintf@libc.so.6":   "8",
	"wscanf@libc.so.6":    "4",
}

func InitDB() error {
	db, err := leveldb.OpenFile(PERMFOLDER, nil)

	if err != nil {
		return err
	}

	defer func() {
		err = db.Close()
	}()

	t, err := db.OpenTransaction()

	if err != nil {
		return err
	}

	for k, v := range perm {
		err = t.Put([]byte(k), []byte(v), nil)

		if err != nil {
			return err
		}
	}

	return t.Commit()
}
