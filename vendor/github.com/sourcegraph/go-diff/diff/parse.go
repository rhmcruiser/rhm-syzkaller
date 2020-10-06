	"path/filepath"
	"strconv"
// ParseMultiFileDiff parses a multi-file unified diff. It returns an error if
// parsing failed as a whole, but does its best to parse as many files in the
// case of per-file errors. If it cannot detect when the diff of the next file
// begins, the hunks are added to the FileDiff of the previous file.
	// FileDiff is added/deleted file
	// No further collection of hunks needed
	if fd.NewName == "" {
		return fd, nil
	}

	if err != nil && err != io.EOF {
		fd.OrigTime = origTime
		fd.NewTime = newTime
// timestamps). Or which starts with "Only in " with dir path and filename.
// "Only in" message is supported in POSIX locale: https://pubs.opengroup.org/onlinepubs/9699919799/utilities/diff.html#tag_20_34_10
	if r.fileHeaderLine != nil {
		if isOnlyMessage, source, filename := parseOnlyInMessage(r.fileHeaderLine); isOnlyMessage {
			return filepath.Join(string(source), string(filename)),
				"", nil, nil, nil
		}
	}

	unquotedOrigName, err := strconv.Unquote(origName)
	if err == nil {
		origName = unquotedOrigName
	}
	unquotedNewName, err := strconv.Unquote(newName)
	if err == nil {
		newName = unquotedNewName
	}

	return fmt.Sprintf("overflowed into next file: %s", string(e))
		// Reached message that file is added/deleted
		if isOnlyInMessage, _, _ := parseOnlyInMessage(line); isOnlyInMessage {
			r.fileHeaderLine = line // pass to readOneFileHeader (see fileHeaderLine field doc)
			return xheaders, nil
		}

	var err error
	lineCount := len(fd.Extended)
	if lineCount > 0 && !strings.HasPrefix(fd.Extended[0], "diff --git ") {
		return false
	}
	case (lineCount == 3 || lineCount == 4 && strings.HasPrefix(fd.Extended[3], "Binary files ") || lineCount > 4 && strings.HasPrefix(fd.Extended[3], "GIT binary patch")) &&
		strings.HasPrefix(fd.Extended[1], "new file mode "):
		fd.NewName, err = strconv.Unquote(names[1])
		if err != nil {
			fd.NewName = names[1]
		}
	case (lineCount == 3 || lineCount == 4 && strings.HasPrefix(fd.Extended[3], "Binary files ") || lineCount > 4 && strings.HasPrefix(fd.Extended[3], "GIT binary patch")) &&
		strings.HasPrefix(fd.Extended[1], "deleted file mode "):
		fd.OrigName, err = strconv.Unquote(names[0])
		if err != nil {
			fd.OrigName = names[0]
		}
	case lineCount == 4 && strings.HasPrefix(fd.Extended[2], "rename from ") && strings.HasPrefix(fd.Extended[3], "rename to "):
		fd.OrigName, err = strconv.Unquote(names[0])
		if err != nil {
			fd.OrigName = names[0]
		}
		fd.NewName, err = strconv.Unquote(names[1])
		if err != nil {
			fd.NewName = names[1]
		}
	case lineCount == 6 && strings.HasPrefix(fd.Extended[5], "Binary files ") && strings.HasPrefix(fd.Extended[2], "rename from ") && strings.HasPrefix(fd.Extended[3], "rename to "):
	case lineCount == 3 && strings.HasPrefix(fd.Extended[2], "Binary files ") || lineCount > 3 && strings.HasPrefix(fd.Extended[2], "GIT binary patch"):
		names := strings.SplitN(fd.Extended[0][len("diff --git "):], " ", 2)
		fd.OrigName, err = strconv.Unquote(names[0])
		if err != nil {
			fd.OrigName = names[0]
		}
		fd.NewName, err = strconv.Unquote(names[1])
		if err != nil {
			fd.NewName = names[1]
		}
		return true

	// ErrBadOnlyInMessage is when a file have a malformed `only in` message
	// Should be in format `Only in {source}: {filename}`
	ErrBadOnlyInMessage = errors.New("bad 'only in' message")
			if len(line) >= 1 && (!linePrefix(line[0]) || bytes.HasPrefix(line, []byte("--- "))) {
// parseOnlyInMessage checks if line is a "Only in {source}: {filename}" and returns source and filename
func parseOnlyInMessage(line []byte) (bool, []byte, []byte) {
	if !bytes.HasPrefix(line, onlyInMessagePrefix) {
		return false, nil, nil
	}
	line = line[len(onlyInMessagePrefix):]
	idx := bytes.Index(line, []byte(": "))
	if idx < 0 {
		return false, nil, nil
	}
	return true, line[:idx], line[idx+2:]
}
