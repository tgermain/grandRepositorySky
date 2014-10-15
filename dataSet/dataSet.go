package dataSet

type DataSet struct {
	set map[string]PieceOfData
}

type PieceOfData struct {
	Tag   string //to which node the data is from
	Value string //the data
}

func (d *DataSet) SetData(key string, value PieceOfData) bool {
	_, exist := d.set[key]
	if exist {
		return false
	} else {
		d.set[key] = value
		return true
	}
}

func (d *DataSet) GetData(key string) PieceOfData {
	//can return nul
	val, _ := d.set[key]
	return val
}

func (d *DataSet) DelData(key string) {
	delete(d.set, key)
}

func (d *DataSet) GetSet() map[string]PieceOfData {
	return d.set
}

func MakeDataSet() DataSet {
	return DataSet{
		make(map[string]PieceOfData),
	}
}
