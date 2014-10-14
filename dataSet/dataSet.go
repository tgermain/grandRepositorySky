package dataSet

type DataSet struct {
	set map[string]PieceOfData
}

type PieceOfData struct {
	tag   string //to which node the data is from
	value string //the data
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

func MakeDataSet() DataSet {
	return DataSet{
		make(map[string]PieceOfData),
	}
}
