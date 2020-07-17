package mergeconf

// MergeConf ...
func MergeConf(src1, src2 []byte) ([]byte, error) {
	conf1, err := initFromBytes(src1)
	if err != nil {
		return nil, err
	}
	conf2, err := initFromBytes(src2)
	if err != nil {
		return nil, err
	}
	merge(conf1, conf2)
	return []byte(conf1.String()), nil
}

func merge(conf1, conf2 *elem) {
	if conf1.name != conf2.name {
		return
	}
	if conf1.kind != conf2.kind {
		return
	}
	if conf1.kind == Leaf {
		conf1.value = conf2.value
		return
	}
	for k := range conf2.children {
		if _, ok := conf1.children[k]; ok {
			merge(conf1.children[k], conf2.children[k])
		} else {
			conf1.children[k] = conf2.children[k]
		}
	}

}
