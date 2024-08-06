package product

import (
	"database/sql"

	"github.com/TheRanomial/GoEcomApi/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]*types.Product,error){

	rows,err:=s.db.Query("SELECT * FROM products")

	if err!=nil{
		return nil,err
	}

	products:=make([]*types.Product,0)

	for rows.Next(){
		p,err:=ScanRowsIntoUser(rows)
		if err!=nil{
			return nil,err
		}
		products = append(products, p)
	}
	return products,nil
}

func (s *Store) GetProductById(id int) (*types.Product,error){

	rows,err:=s.db.Query("SELECT * FROM products WHERE id=?",id)

	if err!=nil{
		return nil,err
	}
	product:=new(types.Product)
	for rows.Next(){
		product,err=ScanRowsIntoUser(rows)
		if err!=nil{
			return nil,err
		}
	}
	return product,nil
}

func (s *Store) CreateProduct(product types.CreateProductPayload) error {
	_, err := s.db.Exec("INSERT INTO products (name, price, image, description, quantity) VALUES (?, ?, ?, ?, ?)", product.Name, product.Price, product.Image, product.Description, product.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func ScanRowsIntoUser(rows *sql.Rows) (*types.Product,error){

	product:=new(types.Product)
	err:=rows.Scan(&product.ID,&product.Name,&product.Description,&product.Image,&product.Price,&product.Quantity)

	if err!=nil{
		return nil,err
	}

	return product,nil
}

