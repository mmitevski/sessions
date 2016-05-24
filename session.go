package sessions

type Session interface {

	Set(key, value interface{}) error //set session value

	Get(key interface{}) interface{}  //get session value

	Delete(key interface{}) error     //delete session value

}

