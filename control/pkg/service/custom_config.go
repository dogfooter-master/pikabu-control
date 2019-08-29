package service

type CustomConfigObject struct {
	Ethnicity CustomList `bson:"ethnicity,omitempty"`
	Country   CustomList `bson:"country,omitempty"`
	Skin      CustomList `bson:"skin,omitempty"`
	Disease   CustomList `bson:"disease,omitempty"`
	Location  CustomList `bson:"location,omitempty"`
	Gender    CustomList `bson:"gender,omitempty"`
	Tag       CustomList `bson:"tag,omitempty"`
}

type CustomList struct {
	Default  string   `bson:"default,omitempty"`
	Added    []string `bson:"added,omitempty"`
	Excluded []string `bson:"excluded,omitempty"`
}

func (c *CustomList) Apply(srcDefault string, srcList []string) (dstDefault string, dstList []string, err error) {
	var cache = make(map[string]bool)
	for _, e := range c.Excluded {
		cache[e] = true
	}
	if cache[c.Default] == true {
		c.Default = ""
	}
	if cache[srcDefault] == true {
		srcDefault = ""
	}
	for _, e := range srcList {
		if cache[e] == false {
			cache[e] = true
			dstList = append(dstList, e)
		}
	}
	for _, e := range c.Added {
		if cache[e] == false {
			cache[e] = true
			dstList = append(dstList, e)
		}
	}
	if len(c.Default) > 0 {
		if cache[c.Default] == false {
			dstList = append(dstList, c.Default)
		}
		dstDefault = c.Default
	} else {
		if len(srcDefault) > 0 {
			dstDefault = srcDefault
		} else {
			if len(dstList) > 0 {
				dstDefault = dstList[0]
			}
		}
	}

	return
}
