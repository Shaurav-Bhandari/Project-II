-- ============================================
-- NEPALI FOOD MENU SEED DATA
-- Pricing referenced from Pathao Food Nepal
-- ============================================
-- WARNING: This script clears existing menu data
-- and non-admin users before re-seeding.
-- ============================================

BEGIN;

-- ============================================
-- 1. CLEAR EXISTING DATA (preserve admin)
-- ============================================

-- Remove order-related data
DELETE FROM kot_items;
DELETE FROM kots;
DELETE FROM payments;
DELETE FROM order_items;
DELETE FROM orders;

-- Remove existing menu data
DELETE FROM menu_items;
DELETE FROM menu_categories;

-- Remove non-admin users
DELETE FROM users WHERE email != 'admin@restaurant.com';

-- ============================================
-- 2. INSERT NEPALI FOOD CATEGORIES
-- ============================================

INSERT INTO menu_categories (name, description, display_order) VALUES
    ('Momo',               'Traditional Nepali dumplings — steamed, fried, jhol & more',  1),
    ('Chowmein & Noodles', 'Stir-fried and soupy noodle dishes',                          2),
    ('Dal Bhat & Thali',   'Complete Nepali meals with rice, dal, and sides',              3),
    ('Thukpa & Soups',     'Hearty Himalayan noodle soups',                                4),
    ('Sekuwa & Grills',    'Grilled and barbecued meats',                                  5),
    ('Beverages',          'Hot and cold drinks',                                          6),
    ('Desserts & Sweets',  'Traditional Nepali sweets and treats',                         7);

-- ============================================
-- 3. INSERT MENU ITEMS
-- ============================================

-- === MOMO ===
INSERT INTO menu_items (category_id, name, description, price, preparation_time, is_vegetarian, is_vegan, spice_level, is_available) VALUES
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Chicken Steam Momo',   '10 pcs — juicy chicken filling, served with tomato achar',        220, 15, false, false, 1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Veg Steam Momo',       '10 pcs — mixed vegetables & tofu filling with sesame achar',      180, 15, true,  true,  0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Buff Steam Momo',      '10 pcs — spiced buffalo momo with jhol achar',                    200, 15, false, false, 1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Chicken Fried Momo',   '10 pcs — crispy pan-fried chicken momo',                          250, 18, false, false, 1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Buff Fried Momo',      '10 pcs — crispy fried buffalo momo with spicy chutney',           250, 18, false, false, 2, true),
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Chicken Jhol Momo',    '10 pcs — chicken momo in tangy sesame-tomato jhol sauce',         280, 20, false, false, 3, true),
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Veg Kothey Momo',      '10 pcs — half steamed, half fried vegetable momo',                220, 18, true,  true,  1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Paneer Momo',          '10 pcs — cottage cheese & spinach filling',                       240, 15, true,  false, 0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Chicken C Momo',       '10 pcs — chili chicken momo tossed in spicy sauce',               300, 20, false, false, 4, true),
    ((SELECT id FROM menu_categories WHERE name = 'Momo'), 'Pork Momo',            '10 pcs — minced pork momo with timur achar',                      260, 15, false, false, 2, true);

-- === CHOWMEIN & NOODLES ===
INSERT INTO menu_items (category_id, name, description, price, preparation_time, is_vegetarian, is_vegan, spice_level, is_available) VALUES
    ((SELECT id FROM menu_categories WHERE name = 'Chowmein & Noodles'), 'Veg Chowmein',      'Stir-fried noodles with seasonal vegetables',          180, 12, true,  true,  1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Chowmein & Noodles'), 'Chicken Chowmein',  'Stir-fried noodles with chicken & vegetables',          250, 15, false, false, 1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Chowmein & Noodles'), 'Buff Chowmein',     'Stir-fried noodles with buffalo meat',                  230, 15, false, false, 1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Chowmein & Noodles'), 'Mixed Chowmein',    'Loaded chowmein with chicken, egg & vegetables',        300, 15, false, false, 2, true),
    ((SELECT id FROM menu_categories WHERE name = 'Chowmein & Noodles'), 'Egg Chowmein',      'Stir-fried noodles with scrambled egg',                  200, 12, true,  false, 1, true);

-- === DAL BHAT & THALI ===
INSERT INTO menu_items (category_id, name, description, price, preparation_time, is_vegetarian, is_vegan, spice_level, is_available) VALUES
    ((SELECT id FROM menu_categories WHERE name = 'Dal Bhat & Thali'), 'Veg Thali Set',     'Rice, dal, 2 tarkari, achar, papad, salad',              350, 10, true,  true,  1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Dal Bhat & Thali'), 'Chicken Thali Set', 'Rice, dal, chicken curry, tarkari, achar, papad, salad',  450, 12, false, false, 2, true),
    ((SELECT id FROM menu_categories WHERE name = 'Dal Bhat & Thali'), 'Mutton Thali Set',  'Rice, dal, mutton curry, tarkari, achar, papad, salad',   550, 15, false, false, 2, true),
    ((SELECT id FROM menu_categories WHERE name = 'Dal Bhat & Thali'), 'Buff Thali Set',    'Rice, dal, buff curry, tarkari, achar, papad, salad',     400, 12, false, false, 2, true),
    ((SELECT id FROM menu_categories WHERE name = 'Dal Bhat & Thali'), 'Fish Thali Set',    'Rice, dal, fried fish, tarkari, achar, papad, salad',     500, 15, false, false, 1, true);

-- === THUKPA & SOUPS ===
INSERT INTO menu_items (category_id, name, description, price, preparation_time, is_vegetarian, is_vegan, spice_level, is_available) VALUES
    ((SELECT id FROM menu_categories WHERE name = 'Thukpa & Soups'), 'Veg Thukpa',      'Hearty Himalayan noodle soup with mixed vegetables',    200, 15, true,  true,  1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Thukpa & Soups'), 'Chicken Thukpa',  'Tibetan noodle soup with shredded chicken',              280, 18, false, false, 1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Thukpa & Soups'), 'Buff Thukpa',     'Rich buffalo noodle soup with herbs',                    250, 18, false, false, 1, true),
    ((SELECT id FROM menu_categories WHERE name = 'Thukpa & Soups'), 'Mushroom Soup',   'Creamy mushroom soup with croutons',                     180, 10, true,  false, 0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Thukpa & Soups'), 'Tomato Soup',     'Classic tomato soup with herbs',                         150, 8,  true,  true,  0, true);

-- === SEKUWA & GRILLS ===
INSERT INTO menu_items (category_id, name, description, price, preparation_time, is_vegetarian, is_vegan, spice_level, is_available) VALUES
    ((SELECT id FROM menu_categories WHERE name = 'Sekuwa & Grills'), 'Chicken Sekuwa',         'Marinated grilled chicken skewers (250g)',       350, 20, false, false, 2, true),
    ((SELECT id FROM menu_categories WHERE name = 'Sekuwa & Grills'), 'Buff Sekuwa',            'Spiced buffalo meat skewers (250g)',              400, 22, false, false, 3, true),
    ((SELECT id FROM menu_categories WHERE name = 'Sekuwa & Grills'), 'Pork Sekuwa',            'Tender grilled pork skewers (250g)',              380, 22, false, false, 2, true),
    ((SELECT id FROM menu_categories WHERE name = 'Sekuwa & Grills'), 'Chicken Tandoori (Half)', 'Half tandoori chicken with mint chutney',        450, 25, false, false, 2, true),
    ((SELECT id FROM menu_categories WHERE name = 'Sekuwa & Grills'), 'Paneer Tikka',           'Grilled cottage cheese with peppers & onions',   320, 18, true,  false, 1, true);

-- === BEVERAGES ===
INSERT INTO menu_items (category_id, name, description, price, preparation_time, is_vegetarian, is_vegan, spice_level, is_available) VALUES
    ((SELECT id FROM menu_categories WHERE name = 'Beverages'), 'Masala Chiya',     'Traditional Nepali spiced milk tea',          60,  3, true, false, 0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Beverages'), 'Lemon Tea',        'Refreshing lemon tea, hot or cold',           50,  3, true, true,  0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Beverages'), 'Lassi (Sweet)',    'Creamy yogurt drink',                        120,  3, true, false, 0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Beverages'), 'Mango Lassi',      'Yogurt blended with mango pulp',             150,  5, true, false, 0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Beverages'), 'Fresh Orange Juice','Freshly squeezed seasonal oranges',          180,  5, true, true,  0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Beverages'), 'Lemon Soda',       'Sparkling lemon soda, sweet or salted',       80,  2, true, true,  0, true);

-- === DESSERTS & SWEETS ===
INSERT INTO menu_items (category_id, name, description, price, preparation_time, is_vegetarian, is_vegan, spice_level, is_available) VALUES
    ((SELECT id FROM menu_categories WHERE name = 'Desserts & Sweets'), 'Rasbari',          'Soft cottage cheese balls in sweet syrup',      150, 3,  true, false, 0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Desserts & Sweets'), 'Jalebi',           'Crispy deep-fried swirls soaked in saffron syrup', 120, 5, true, true,  0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Desserts & Sweets'), 'Kheer',            'Creamy rice pudding with cardamom & nuts',      180, 5,  true, false, 0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Desserts & Sweets'), 'Sikarni',          'Sweetened strained yogurt with dry fruits',      160, 3,  true, false, 0, true),
    ((SELECT id FROM menu_categories WHERE name = 'Desserts & Sweets'), 'Sel Roti',         'Traditional Nepali ring-shaped rice bread',      100, 8,  true, true,  0, true);

-- ============================================
-- 4. INSERT STAFF & KITCHEN USERS
-- ============================================
-- Password for all seeded users: staff123
-- Bcrypt hash: $2a$10$7NL.6BBonR0k7qC3/NRIlu2J0HZhosTceKLe10MZ35ZsvDMliqP/C

-- Staff (waiters)
INSERT INTO users (email, password_hash, first_name, last_name, role_id, is_active) VALUES
    ('waiter1@restaurant.com', '$2a$10$7NL.6BBonR0k7qC3/NRIlu2J0HZhosTceKLe10MZ35ZsvDMliqP/C', 'Ram',    'Thapa',   (SELECT id FROM roles WHERE name = 'staff'),   true),
    ('waiter2@restaurant.com', '$2a$10$7NL.6BBonR0k7qC3/NRIlu2J0HZhosTceKLe10MZ35ZsvDMliqP/C', 'Sita',   'Gurung',  (SELECT id FROM roles WHERE name = 'staff'),   true),
    ('waiter3@restaurant.com', '$2a$10$7NL.6BBonR0k7qC3/NRIlu2J0HZhosTceKLe10MZ35ZsvDMliqP/C', 'Hari',   'Tamang',  (SELECT id FROM roles WHERE name = 'staff'),   true);

-- Kitchen staff
INSERT INTO users (email, password_hash, first_name, last_name, role_id, is_active) VALUES
    ('kitchen1@restaurant.com', '$2a$10$7NL.6BBonR0k7qC3/NRIlu2J0HZhosTceKLe10MZ35ZsvDMliqP/C', 'Bishnu', 'Sherpa', (SELECT id FROM roles WHERE name = 'kitchen'), true),
    ('kitchen2@restaurant.com', '$2a$10$7NL.6BBonR0k7qC3/NRIlu2J0HZhosTceKLe10MZ35ZsvDMliqP/C', 'Maya',   'Rai',    (SELECT id FROM roles WHERE name = 'kitchen'), true);

-- Manager
INSERT INTO users (email, password_hash, first_name, last_name, role_id, is_active) VALUES
    ('manager@restaurant.com', '$2a$10$7NL.6BBonR0k7qC3/NRIlu2J0HZhosTceKLe10MZ35ZsvDMliqP/C', 'Prakash', 'Shrestha', (SELECT id FROM roles WHERE name = 'manager'), true);

COMMIT;

-- ============================================
-- SEEDED CREDENTIALS SUMMARY
-- ============================================
-- | Email                      | Password  | Role    |
-- |--------------------------- |-----------|---------|
-- | admin@restaurant.com       | admin123  | admin   |
-- | manager@restaurant.com     | staff123  | manager |
-- | waiter1@restaurant.com     | staff123  | staff   |
-- | waiter2@restaurant.com     | staff123  | staff   |
-- | waiter3@restaurant.com     | staff123  | staff   |
-- | kitchen1@restaurant.com    | staff123  | kitchen |
-- | kitchen2@restaurant.com    | staff123  | kitchen |
-- ============================================
