package types

import (
	"errors"
	"strconv"

	"github.com/YendisFish/sirius/toybox/helpers"
)

type PasswdUser struct {
	Username string
	Haspassword byte
	UID int
	GID int
	UIDInf string
	Home string
	Shell string
}

func ReadPasswd() ([]PasswdUser, error) {
	fle, err := helpers.ReadFileLns("/etc/passwd")
	if err != nil {
		return nil, err
	}

	var ret []PasswdUser
	for _, ln := range fle {
		if len(ln) < 14 {
			continue
		}

		parsed := helpers.SliceUp(ln, ':')

		var entry PasswdUser
		err := passwdEntryFromSlice(parsed, &entry)

		if err != nil {
			return nil, err
		}

		ret = append(ret, entry)
	}

	return ret, nil
}

func passwdEntryFromSlice(buff [][]byte, out *PasswdUser) error {
	if len(buff) < 7 {
		return errors.New("Invalid passwd entry length")
	}

	uname := string(buff[0])
	haspass := buff[1][0]

	userId, err := strconv.Atoi(string(buff[2]))
	if err != nil {
		return err
	}

	groupId, err := strconv.Atoi(string(buff[3]))
	if err != nil {
		return err
	}

	userIdInf := string(buff[4])
	home := string(buff[5])
	shell := string(buff[6])

	*out = PasswdUser{
		Username: uname,
		Haspassword: haspass,
		UID: userId,
		GID: groupId,
		UIDInf: userIdInf,
		Home: home,
		Shell: shell,
	}

	return nil
}
