-- Script de inicialización de la base de datos
-- Este archivo se ejecuta automáticamente cuando se crea el contenedor MySQL

-- Crear base de datos si no existe
CREATE DATABASE IF NOT EXISTS padel_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Usar la base de datos
USE padel_db;

-- Crear usuario si no existe (opcional, ya se crea con variables de entorno)
-- CREATE USER IF NOT EXISTS 'padel_user'@'%' IDENTIFIED BY 'padel_password';
-- GRANT ALL PRIVILEGES ON padel_db.* TO 'padel_user'@'%';
-- FLUSH PRIVILEGES;

-- Las tablas se crean automáticamente mediante GORM AutoMigrate
-- Este archivo se puede usar para datos iniciales o configuraciones adicionales
