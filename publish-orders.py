import requests
import uuid
import random

# Generate 1000 unique item IDs using UUIDs
item_ids = [uuid.uuid4().__str__() for _ in range(1000)]

# Generate 100 unique customer IDs using UUIDs
customers = [uuid.uuid4().__str__() for _ in range(100)]

# Generate and post 120 orders
for i in range(120):
    # Randomly select a customer ID for the order
    customer = random.choice(customers)

    # Randomly determine the number of line items for the order (between 1 and 10)
    num_line_items = random.randint(1, 10)

    # Generate line items with random item IDs, quantities, and prices
    line_items = [
        {
            "item_id": random.choice(item_ids),
            "quantity": random.randint(1, 10),
            "price": random.randint(1, 10000),
        }
        for _ in range(num_line_items)
    ]

    # Create the order with customer ID and line items
    order = {
        "customer_id": customer,
        "line_items": line_items,
    }

    # Post the order to the specified URL (http://localhost:3000/orders)
    r = requests.post("http://localhost:3000/orders", json=order)

    # Accessing the status code of the HTTP response (but not doing anything with it)
    r.status_code
    
    # Print a message indicating the order has been posted
    print("posted order", i + 1)