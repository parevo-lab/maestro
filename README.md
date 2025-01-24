# Go Workflow Management

## About

This project demonstrates a robust workflow management system implementation in Go. It provides a flexible and scalable solution for managing complex business processes through a RESTful API. The system supports various workflow types with customizable steps, including tasks, approvals, and automated processes.

### Key Features
- 🔄 Dynamic workflow creation and management
- 👥 User assignment and role-based operations
- ✅ Multi-step approval processes
- 📊 Data validation at each step
- 🔍 Process tracking and monitoring
- 📝 Detailed documentation with examples
- 🔔 Real-time notifications via WebSocket
- 💾 Persistent workflow state management

## API Endpoints

### Workflow Management
- `POST /workflows` - Create a new workflow
- `GET /workflows/{id}` - Get workflow details by ID
- `GET /workflows/user/{userId}` - Get all workflows assigned to a user
- `POST /workflows/{id}/steps/{stepId}/process` - Process a workflow step

### WebSocket Notifications
- `GET /ws` - WebSocket connection endpoint for real-time notifications

## Workflow Types
- `task` - Basic task step
- `approval` - Approval required step
- `decision` - Decision making step
- `process` - Automated process step

## Result Types
- `invoice` - Invoice generation result
- `document` - Document generation result
- `report` - Report generation result
- `notification` - Notification result

## Notification Types
- `workflow` - Workflow related notifications
- `task` - Task related notifications
- `system` - System notifications

## Example Usage

### 1. Creating Order Workflow

```bash
curl -X POST http://localhost:8080/workflows \
-H "Content-Type: application/json" \
-d '{
  "name": "Order Process",
  "type": "order_process",
  "created_by": "65b012345678901234567890",
  "steps": [
    {
      "id": "65b012345678901234567891",
      "type": "task",
      "title": "Order Details",
      "assigned_to": "65b012345678901234567890",
      "status": "pending",
      "next_steps": ["65b012345678901234567892"]
    },
    {
      "id": "65b012345678901234567892",
      "type": "approval",
      "title": "Stock Control",
      "assigned_to": "65b012345678901234567893",
      "status": "pending",
      "next_steps": ["65b012345678901234567894"],
      "required_data": ["order_items", "total_amount"]
    },
    {
      "id": "65b012345678901234567894",
      "type": "process",
      "title": "Invoice Generation",
      "assigned_to": "65b012345678901234567895",
      "status": "pending",
      "result_type": "invoice",
      "required_data": ["order_items", "customer_info", "total_amount", "stock_approval"]
    }
  ]
}'
```

### 2. Entering Order Details

```bash
curl -X POST http://localhost:8080/workflows/WORKFLOW_ID/steps/65b012345678901234567891/process \
-H "Content-Type: application/json" \
-d '{
  "action": "approve",
  "data": {
    "order_items": [
      {
        "product_id": "PROD001",
        "name": "Laptop",
        "quantity": 1,
        "price": 15000
      }
    ],
    "customer_info": {
      "name": "John Smith",
      "email": "john@example.com",
      "tax_number": "1234567890"
    },
    "total_amount": 15000
  }
}'
```

### 3. Stock Control

```bash
curl -X POST http://localhost:8080/workflows/WORKFLOW_ID/steps/65b012345678901234567892/process \
-H "Content-Type: application/json" \
-d '{
  "action": "approve",
  "data": {
    "stock_approval": true,
    "stock_notes": "Stock is sufficient",
    "approved_by": "Warehouse Manager"
  }
}'
```

### 4. Invoice Generation

```bash
curl -X POST http://localhost:8080/workflows/WORKFLOW_ID/steps/65b012345678901234567894/process \
-H "Content-Type: application/json" \
-d '{
  "action": "approve",
  "data": {
    "invoice_number": "INV-2024-001",
    "invoice_date": "2024-01-24T15:00:00Z",
    "items": [
      {
        "product_id": "PROD001",
        "name": "Laptop",
        "quantity": 1,
        "unit_price": 15000,
        "total": 15000
      }
    ],
    "subtotal": 15000,
    "tax": 2700,
    "total": 17700
  }
}'
```
