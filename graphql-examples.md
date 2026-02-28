# GraphQL Playground Examples

Use these examples in the GraphQL Playground at `http://localhost:8000/playground`

## Step-by-Step Test Flow

### 1. Create Accounts

```graphql
mutation CreateAccount1 {
  createAccount(account: { name: "John Doe" }) {
    id
    name
  }
}
```

```graphql
mutation CreateAccount2 {
  createAccount(account: { name: "Jane Smith" }) {
    id
    name
  }
}
```

```graphql
mutation CreateAccount3 {
  createAccount(account: { name: "Bob Johnson" }) {
    id
    name
  }
}
```

**Expected Response:**
```json
{
  "data": {
    "createAccount": {
      "id": "01HXXXXX...",
      "name": "John Doe"
    }
  }
}
```

**Note:** Save the `id` values from these mutations - you'll need them for creating orders!

---

### 2. Create Products

```graphql
mutation CreateProduct1 {
  createProduct(product: {
    name: "Laptop"
    description: "High-performance laptop with 16GB RAM"
    price: 1299.99
  }) {
    id
    name
    description
    price
  }
}
```

```graphql
mutation CreateProduct2 {
  createProduct(product: {
    name: "Wireless Mouse"
    description: "Ergonomic wireless mouse with long battery life"
    price: 29.99
  }) {
    id
    name
    description
    price
  }
}
```

```graphql
mutation CreateProduct3 {
  createProduct(product: {
    name: "Mechanical Keyboard"
    description: "RGB mechanical keyboard with Cherry MX switches"
    price: 149.99
  }) {
    id
    name
    description
    price
  }
}
```

```graphql
mutation CreateProduct4 {
  createProduct(product: {
    name: "Monitor"
    description: "27-inch 4K IPS monitor"
    price: 399.99
  }) {
    id
    name
    description
    price
  }
}
```

**Note:** Save the product `id` values - you'll need them for creating orders!

---

### 3. Query All Accounts

```graphql
query GetAllAccounts {
  accounts {
    id
    name
  }
}
```

---

### 4. Query Accounts with Pagination

```graphql
query GetAccountsPaginated {
  accounts(paginationInput: { skip: 0, take: 10 }) {
    id
    name
  }
}
```

```graphql
query GetAccountsPaginated2 {
  accounts(paginationInput: { skip: 1, take: 2 }) {
    id
    name
  }
}
```

---

### 5. Query Account by ID

```graphql
query GetAccountById {
  accounts(id: "PASTE_ACCOUNT_ID_HERE") {
    id
    name
    orders {
      id
      createdAt
      totalPrice
      products {
        id
        name
        description
        price
        quantity
      }
    }
  }
}
```

**Replace `PASTE_ACCOUNT_ID_HERE` with an actual account ID from step 1**

---

### 6. Query All Products

```graphql
query GetAllProducts {
  products {
    id
    name
    description
    price
  }
}
```

---

### 7. Query Products with Pagination

```graphql
query GetProductsPaginated {
  products(paginationInput: { skip: 0, take: 5 }) {
    id
    name
    description
    price
  }
}
```

```graphql
query GetProductsPaginated2 {
  products(paginationInput: { skip: 2, take: 2 }) {
    id
    name
    description
    price
  }
}
```

---

### 8. Search Products by Query

```graphql
query SearchProducts {
  products(query: "laptop", paginationInput: { skip: 0, take: 10 }) {
    id
    name
    description
    price
  }
}
```

```graphql
query SearchProducts2 {
  products(query: "keyboard", paginationInput: { skip: 0, take: 10 }) {
    id
    name
    description
    price
  }
}
```

```graphql
query SearchProducts3 {
  products(query: "wireless", paginationInput: { skip: 0, take: 10 }) {
    id
    name
    description
    price
  }
}
```

---

### 9. Query Product by ID

```graphql
query GetProductById {
  products(id: "PASTE_PRODUCT_ID_HERE") {
    id
    name
    description
    price
  }
}
```

**Replace `PASTE_PRODUCT_ID_HERE` with an actual product ID from step 2**

---

### 10. Create Orders

**Important:** Replace `ACCOUNT_ID` and `PRODUCT_ID` values with actual IDs from previous mutations!

```graphql
mutation CreateOrder1 {
  createOrder(order: {
    accountId: "PASTE_ACCOUNT_ID_HERE"
    products: [
      {
        id: "PASTE_PRODUCT_ID_1_HERE"
        quantity: 1
      }
      {
        id: "PASTE_PRODUCT_ID_2_HERE"
        quantity: 2
      }
    ]
  }) {
    id
    createdAt
    totalPrice
    products {
      id
      name
      description
      price
      quantity
    }
  }
}
```

**Example with actual values:**
```graphql
mutation CreateOrderExample {
  createOrder(order: {
    accountId: "01HXXXXX..."
    products: [
      {
        id: "01HYYYYY..."
        quantity: 1
      }
      {
        id: "01HZZZZZ..."
        quantity: 3
      }
    ]
  }) {
    id
    createdAt
    totalPrice
    products {
      id
      name
      description
      price
      quantity
    }
  }
}
```

---

### 11. Query Account with Orders (Full Relationship)

```graphql
query GetAccountWithOrders {
  accounts(id: "PASTE_ACCOUNT_ID_HERE") {
    id
    name
    orders {
      id
      createdAt
      totalPrice
      products {
        id
        name
        description
        price
        quantity
      }
    }
  }
}
```

---

## Complete Test Scenario

Run these in order to test the full flow:

1. **Create 2 accounts** (use mutations from step 1)
2. **Create 3 products** (use mutations from step 2)
3. **Query all accounts** (step 3)
4. **Query all products** (step 6)
5. **Create an order** (step 10) - Use IDs from steps 1 & 2
6. **Query account with orders** (step 11) - See the order relationship

---

## Error Testing

### Invalid Order (zero quantity)
```graphql
mutation CreateOrderInvalid {
  createOrder(order: {
    accountId: "PASTE_ACCOUNT_ID_HERE"
    products: [
      {
        id: "PASTE_PRODUCT_ID_HERE"
        quantity: 0
      }
    ]
  }) {
    id
  }
}
```

**Expected:** Error message about invalid quantity

### Invalid Order (negative quantity)
```graphql
mutation CreateOrderInvalid2 {
  createOrder(order: {
    accountId: "PASTE_ACCOUNT_ID_HERE"
    products: [
      {
        id: "PASTE_PRODUCT_ID_HERE"
        quantity: -1
      }
    ]
  }) {
    id
  }
}
```

---

## Tips

1. **Copy IDs:** After each mutation, copy the returned `id` values - you'll need them for subsequent queries/mutations
2. **Variables:** You can use GraphQL variables for cleaner queries:
   ```graphql
   query GetAccount($accountId: String!) {
     accounts(id: $accountId) {
       id
       name
     }
   }
   ```
   Variables:
   ```json
   {
     "accountId": "01HXXXXX..."
   }
   ```

3. **Multiple Operations:** You can run multiple queries in one request:
   ```graphql
   query MultipleQueries {
     accounts {
       id
       name
     }
     products {
       id
       name
       price
     }
   }
   ```
