CREATE TABLE supermarket_sales (
	"Invoice" VARCHAR NOT NULL, 
	"StockCode" VARCHAR NOT NULL, 
	"Description" VARCHAR, 
	"Quantity" DECIMAL NOT NULL, 
	"InvoiceDate" DATE NOT NULL, 
	"InvoiceTime" INTERVAL NOT NULL, 
	"Price" DECIMAL NOT NULL, 
	"Customer ID" DECIMAL, 
	"Country" VARCHAR NOT NULL
);
