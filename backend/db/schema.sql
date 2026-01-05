-- Restaurant Operations Management System
-- PostgreSQL Database Schema with RLS

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================
-- ROLES AND USERS
-- ============================================

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    permissions JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role_id UUID REFERENCES roles(id),
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- RESTAURANT TABLES
-- ============================================

CREATE TABLE restaurant_tables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    table_number VARCHAR(20) UNIQUE NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 4,
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'occupied', 'reserved', 'maintenance')),
    location VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- MENU MANAGEMENT
-- ============================================

CREATE TABLE menu_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    display_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE menu_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID REFERENCES menu_categories(id) ON DELETE SET NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    preparation_time INTEGER DEFAULT 15, -- in minutes
    is_available BOOLEAN DEFAULT true,
    is_vegetarian BOOLEAN DEFAULT false,
    is_vegan BOOLEAN DEFAULT false,
    is_gluten_free BOOLEAN DEFAULT false,
    spice_level INTEGER DEFAULT 0 CHECK (spice_level >= 0 AND spice_level <= 5),
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for binary search on menu items by ID
CREATE INDEX idx_menu_items_id ON menu_items(id);
CREATE INDEX idx_menu_items_category ON menu_items(category_id);

-- ============================================
-- ORDERS
-- ============================================

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_number SERIAL,
    table_id UUID REFERENCES restaurant_tables(id),
    customer_name VARCHAR(200),
    status VARCHAR(30) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'preparing', 'ready', 'served', 'completed', 'cancelled')),
    order_type VARCHAR(20) DEFAULT 'dine-in' CHECK (order_type IN ('dine-in', 'takeaway', 'delivery')),
    notes TEXT,
    subtotal DECIMAL(10, 2) DEFAULT 0,
    tax DECIMAL(10, 2) DEFAULT 0,
    discount DECIMAL(10, 2) DEFAULT 0,
    total DECIMAL(10, 2) DEFAULT 0,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id UUID REFERENCES menu_items(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10, 2) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    special_instructions TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'preparing', 'ready', 'served', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for order queries
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_table ON orders(table_id);
CREATE INDEX idx_order_items_order ON order_items(order_id);

-- ============================================
-- KITCHEN ORDER TICKETS (KOT)
-- ============================================

CREATE TABLE kots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kot_number SERIAL,
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled')),
    priority INTEGER DEFAULT 0, -- Higher number = higher priority
    station VARCHAR(50), -- e.g., 'grill', 'salad', 'dessert'
    assigned_chef VARCHAR(100),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE kot_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kot_id UUID REFERENCES kots(id) ON DELETE CASCADE,
    order_item_id UUID REFERENCES order_items(id),
    menu_item_name VARCHAR(200) NOT NULL,
    quantity INTEGER NOT NULL,
    special_instructions TEXT,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for KOT queries (for merge sort prioritization)
CREATE INDEX idx_kots_status ON kots(status);
CREATE INDEX idx_kots_priority_created ON kots(priority DESC, created_at ASC);
CREATE INDEX idx_kot_items_kot ON kot_items(kot_id);

-- ============================================
-- BILLING AND PAYMENTS
-- ============================================

CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID REFERENCES orders(id),
    payment_method_id UUID REFERENCES payment_methods(id),
    amount DECIMAL(10, 2) NOT NULL,
    tip DECIMAL(10, 2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    transaction_reference VARCHAR(100),
    processed_by UUID REFERENCES users(id),
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for payment queries
CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);

-- ============================================
-- SESSION TOKENS (for secure token generation)
-- ============================================

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

-- ============================================
-- ROW LEVEL SECURITY (RLS)
-- ============================================

-- Enable RLS on tables
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE kots ENABLE ROW LEVEL SECURITY;
ALTER TABLE payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Create app roles
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'app_admin') THEN
        CREATE ROLE app_admin;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'app_staff') THEN
        CREATE ROLE app_staff;
    END IF;
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'app_kitchen') THEN
        CREATE ROLE app_kitchen;
    END IF;
END
$$;

-- Admin can access everything
CREATE POLICY admin_all_orders ON orders FOR ALL TO app_admin USING (true);
CREATE POLICY admin_all_order_items ON order_items FOR ALL TO app_admin USING (true);
CREATE POLICY admin_all_kots ON kots FOR ALL TO app_admin USING (true);
CREATE POLICY admin_all_payments ON payments FOR ALL TO app_admin USING (true);
CREATE POLICY admin_all_users ON users FOR ALL TO app_admin USING (true);

-- Staff can manage orders they created or all active orders
CREATE POLICY staff_orders ON orders FOR ALL TO app_staff 
    USING (created_by = current_setting('app.current_user_id', true)::uuid 
           OR status NOT IN ('completed', 'cancelled'));

CREATE POLICY staff_order_items ON order_items FOR ALL TO app_staff 
    USING (EXISTS (SELECT 1 FROM orders o WHERE o.id = order_id 
                   AND (o.created_by = current_setting('app.current_user_id', true)::uuid 
                        OR o.status NOT IN ('completed', 'cancelled'))));

-- Kitchen staff can view and update KOTs
CREATE POLICY kitchen_view_kots ON kots FOR SELECT TO app_kitchen USING (true);
CREATE POLICY kitchen_update_kots ON kots FOR UPDATE TO app_kitchen USING (true);

-- Staff can view their own user record
CREATE POLICY staff_own_user ON users FOR SELECT TO app_staff 
    USING (id = current_setting('app.current_user_id', true)::uuid);

-- ============================================
-- INITIAL DATA
-- ============================================

-- Insert default roles
INSERT INTO roles (name, description, permissions) VALUES
    ('admin', 'System Administrator', '{"all": true}'),
    ('manager', 'Restaurant Manager', '{"orders": true, "menu": true, "reports": true, "users": false}'),
    ('staff', 'Restaurant Staff', '{"orders": true, "menu": false, "reports": false}'),
    ('kitchen', 'Kitchen Staff', '{"kots": true}');

-- Insert default payment methods
INSERT INTO payment_methods (name) VALUES
    ('Cash'),
    ('Credit Card'),
    ('Debit Card'),
    ('Mobile Payment'),
    ('Gift Card');

-- Insert sample tables
INSERT INTO restaurant_tables (table_number, capacity, location) VALUES
    ('T1', 2, 'Window'),
    ('T2', 2, 'Window'),
    ('T3', 4, 'Center'),
    ('T4', 4, 'Center'),
    ('T5', 6, 'Corner'),
    ('T6', 6, 'Corner'),
    ('T7', 8, 'Private'),
    ('T8', 4, 'Patio'),
    ('T9', 4, 'Patio'),
    ('T10', 2, 'Bar');

-- Insert sample menu categories
INSERT INTO menu_categories (name, description, display_order) VALUES
    ('Appetizers', 'Start your meal right', 1),
    ('Main Courses', 'Delicious entrees', 2),
    ('Desserts', 'Sweet endings', 3),
    ('Beverages', 'Drinks and refreshments', 4),
    ('Specials', 'Chef''s special items', 5);

-- Insert sample menu items
INSERT INTO menu_items (category_id, name, description, price, preparation_time, is_vegetarian, spice_level) VALUES
    ((SELECT id FROM menu_categories WHERE name = 'Appetizers'), 'Spring Rolls', 'Crispy vegetable spring rolls with sweet chili sauce', 8.99, 10, true, 1),
    ((SELECT id FROM menu_categories WHERE name = 'Appetizers'), 'Chicken Wings', 'Spicy buffalo wings with blue cheese dip', 12.99, 15, false, 3),
    ((SELECT id FROM menu_categories WHERE name = 'Appetizers'), 'Soup of the Day', 'Ask your server for today''s selection', 6.99, 5, true, 0),
    ((SELECT id FROM menu_categories WHERE name = 'Main Courses'), 'Grilled Salmon', 'Atlantic salmon with lemon butter sauce', 24.99, 20, false, 0),
    ((SELECT id FROM menu_categories WHERE name = 'Main Courses'), 'Ribeye Steak', '12oz prime cut with seasonal vegetables', 34.99, 25, false, 0),
    ((SELECT id FROM menu_categories WHERE name = 'Main Courses'), 'Vegetable Pasta', 'Penne with roasted vegetables in marinara sauce', 16.99, 15, true, 1),
    ((SELECT id FROM menu_categories WHERE name = 'Main Courses'), 'Chicken Tikka Masala', 'Tender chicken in creamy tomato curry', 18.99, 20, false, 2),
    ((SELECT id FROM menu_categories WHERE name = 'Desserts'), 'Chocolate Lava Cake', 'Warm chocolate cake with molten center', 9.99, 12, true, 0),
    ((SELECT id FROM menu_categories WHERE name = 'Desserts'), 'Cheesecake', 'New York style with berry compote', 8.99, 5, true, 0),
    ((SELECT id FROM menu_categories WHERE name = 'Beverages'), 'Fresh Lemonade', 'Homemade with fresh lemons', 4.99, 2, true, 0),
    ((SELECT id FROM menu_categories WHERE name = 'Beverages'), 'Iced Tea', 'Choice of regular or peach', 3.99, 2, true, 0),
    ((SELECT id FROM menu_categories WHERE name = 'Beverages'), 'Coffee', 'Freshly brewed', 3.49, 2, true, 0);

-- Create default admin user (password: admin123)
-- Password hash for 'admin123' using bcrypt
INSERT INTO users (email, password_hash, first_name, last_name, role_id) VALUES
    ('admin@restaurant.com', '$2a$10$eEKnAa13H6n5TYqEP9wSJOgxI5QtcY9sFLLsLWvBJPWrFe5VUg8nK', 'System', 'Admin', (SELECT id FROM roles WHERE name = 'admin'));

-- ============================================
-- FUNCTIONS AND TRIGGERS
-- ============================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at trigger to all relevant tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tables_updated_at BEFORE UPDATE ON restaurant_tables FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON menu_categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_items_updated_at BEFORE UPDATE ON menu_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_kots_updated_at BEFORE UPDATE ON kots FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to calculate order totals
CREATE OR REPLACE FUNCTION calculate_order_total()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE orders SET
        subtotal = (SELECT COALESCE(SUM(total_price), 0) FROM order_items WHERE order_id = NEW.order_id),
        tax = (SELECT COALESCE(SUM(total_price), 0) * 0.10 FROM order_items WHERE order_id = NEW.order_id),
        total = (SELECT COALESCE(SUM(total_price), 0) * 1.10 FROM order_items WHERE order_id = NEW.order_id)
    WHERE id = NEW.order_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER calculate_order_total_trigger 
AFTER INSERT OR UPDATE OR DELETE ON order_items 
FOR EACH ROW EXECUTE FUNCTION calculate_order_total();

-- Function to auto-generate KOT on order confirmation
CREATE OR REPLACE FUNCTION generate_kot_on_confirmation()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'confirmed' AND OLD.status = 'pending' THEN
        INSERT INTO kots (order_id, priority)
        VALUES (NEW.id, 0);
        
        -- Copy order items to KOT items
        INSERT INTO kot_items (kot_id, order_item_id, menu_item_name, quantity, special_instructions)
        SELECT 
            (SELECT id FROM kots WHERE order_id = NEW.id ORDER BY created_at DESC LIMIT 1),
            oi.id,
            mi.name,
            oi.quantity,
            oi.special_instructions
        FROM order_items oi
        JOIN menu_items mi ON oi.menu_item_id = mi.id
        WHERE oi.order_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER generate_kot_trigger
AFTER UPDATE ON orders
FOR EACH ROW EXECUTE FUNCTION generate_kot_on_confirmation();

-- Function to update table status when order is created/completed
CREATE OR REPLACE FUNCTION update_table_status()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.table_id IS NOT NULL THEN
        UPDATE restaurant_tables SET status = 'occupied' WHERE id = NEW.table_id;
    ELSIF TG_OP = 'UPDATE' AND NEW.status IN ('completed', 'cancelled') AND NEW.table_id IS NOT NULL THEN
        UPDATE restaurant_tables SET status = 'available' WHERE id = NEW.table_id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_table_status_trigger
AFTER INSERT OR UPDATE ON orders
FOR EACH ROW EXECUTE FUNCTION update_table_status();
