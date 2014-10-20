package dataSet

const DATANOTFOUND = "no data with this key"

type DataSet struct {
	Set map[string]PieceOfData
}

type PieceOfData struct {
	Tag   string //to which node the data is from
	Value string //the data
}

func (d *DataSet) SetData(key, value, tag string) bool {
	_, exist := d.Set[key]
	if exist {
		return false
	} else {
		d.Set[key] = PieceOfData{
			tag,
			value,
		}
		return true
	}
}

func (d *DataSet) GetData(key string) PieceOfData {
	//can return nul
	val, ok := d.Set[key]
	if ok {
		return val

	} else {
		return PieceOfData{
			"",
			DATANOTFOUND,
		}
	}
}

func (d *DataSet) DelData(key string) {
	delete(d.Set, key)
}

func (d *DataSet) GetSet() map[string]PieceOfData {
	return d.Set
}

func MakeDataSet() DataSet {
	return DataSet{
		make(map[string]PieceOfData),
	}
}
