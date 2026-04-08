# TransactionManager

**cURL Examples**

Base URL: `http://localhost:8080`

Create account:
```bash
curl -X POST http://localhost:8080/accounts \
  -H 'Content-Type: application/json' \
  -d '{"document_number":"12345678900"}'
```

Get account:
```bash
curl http://localhost:8080/accounts/1
```

Create transaction:
Notes:
`amount` must be positive and have up to 2 decimals. The signed amount is derived from `operation_type_id` (`debit` => negative, `credit` => positive).

```bash
curl -X POST http://localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"account_id":1,"operation_type_id":4,"amount":123.45}'
```
