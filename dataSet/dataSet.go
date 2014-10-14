package dataSet

type DataSet struct {
	set map[string]pieceOfData
}

type pieceOfData struct {
	tag   string //to which node the data is from
	value string //the data
}

func (d *DataSet) SetData(key string, value pieceOfData) bool {
	_, exist := d.set[key]
	if exist {
		return false
	} else {
		d.set[key] = value
		return true
	}
}

func (d *DataSet) GetData(key string) pieceOfData {
	//can return nul
	val, _ := d.set[key]
	return val
}

func MakeDataSet() DataSet {
	return DataSet{
		make(map[string]pieceOfData),
	}
}
