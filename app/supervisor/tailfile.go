package supervisor

import (
	"os"
)

func tailFile(path string, maxSize int64) (int64, []byte) {
	st, _ := os.Stat(path)
	if st == nil || st.Size() == 0 {
		return 0, nil
	}
	ff, err := os.Open(path)
	if err != nil {
		return 0, nil
	}
	defer ff.Close()
	var bs []byte
	var offset int64 = 0
	if maxSize < st.Size() {
		bs = make([]byte, int(maxSize))
		offset = st.Size() - maxSize
		ff.ReadAt(bs, offset)
	} else {
		bs = make([]byte, int(st.Size()))
		ff.Read(bs)
	}
	return offset, bs
}
