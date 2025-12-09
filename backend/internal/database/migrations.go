package database

import (
	"log"

	"backend/internal/utils"

	"github.com/google/uuid"
)

func RunMigrations() error {
	log.Println("Running database migrations...")

	queries := []string{
		// Extensión para UUID
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,

		// Tabla de workspaces
		`CREATE TABLE IF NOT EXISTS workspaces (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		// Tabla de usuarios
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255),
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL CHECK (role IN ('super_admin', 'admin', 'client')),
			session_token VARCHAR(255),
			workspace_id UUID REFERENCES workspaces(id) ON DELETE SET NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CHECK (
				(role = 'client' AND workspace_id IS NOT NULL) OR
				(role IN ('admin', 'super_admin') AND workspace_id IS NULL)
			)
		);`,

		// Tabla de relación workspace-admins (muchos a muchos)
		`CREATE TABLE IF NOT EXISTS workspace_admins (
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (workspace_id, admin_id)
		);`,

		// Tabla de grupos de admins
		`CREATE TABLE IF NOT EXISTS admin_groups (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			admin_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		// Tabla de chats
		`CREATE TABLE IF NOT EXISTS chats (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			client_id UUID REFERENCES users(id) ON DELETE CASCADE,
			assigned_to_id UUID REFERENCES users(id) ON DELETE SET NULL,
			is_admin_chat BOOLEAN NOT NULL DEFAULT FALSE,
			other_admin_id UUID REFERENCES users(id) ON DELETE CASCADE,
			workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			last_message_at TIMESTAMP,
			CHECK (
				(is_admin_chat = TRUE AND other_admin_id IS NOT NULL AND client_id IS NULL) OR
				(is_admin_chat = FALSE AND client_id IS NOT NULL AND other_admin_id IS NULL)
			),
			CHECK (
				(is_admin_chat = FALSE AND workspace_id IS NOT NULL) OR
				(is_admin_chat = TRUE)
			)
		);`,

		// Tabla de mensajes
		`CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
			sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT,
			audio_data BYTEA,
			image_data BYTEA,
			mime_type VARCHAR(100),
			is_audio BOOLEAN NOT NULL DEFAULT FALSE,
			is_image BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CHECK (
				(content IS NOT NULL AND audio_data IS NULL AND image_data IS NULL) OR
				(is_audio = TRUE AND audio_data IS NOT NULL) OR
				(is_image = TRUE AND image_data IS NOT NULL)
			)
		);`,

		// Índices
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`,
		`CREATE INDEX IF NOT EXISTS idx_users_session_token ON users(session_token);`,
		`CREATE INDEX IF NOT EXISTS idx_users_workspace_id ON users(workspace_id);`,
		`CREATE INDEX IF NOT EXISTS idx_workspaces_slug ON workspaces(slug);`,
		`CREATE INDEX IF NOT EXISTS idx_workspace_admins_workspace_id ON workspace_admins(workspace_id);`,
		`CREATE INDEX IF NOT EXISTS idx_workspace_admins_admin_id ON workspace_admins(admin_id);`,
		`CREATE INDEX IF NOT EXISTS idx_chats_admin_id ON chats(admin_id);`,
		`CREATE INDEX IF NOT EXISTS idx_chats_client_id ON chats(client_id);`,
		`CREATE INDEX IF NOT EXISTS idx_chats_assigned_to_id ON chats(assigned_to_id);`,
		`CREATE INDEX IF NOT EXISTS idx_chats_workspace_id ON chats(workspace_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);`,
		`CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_chats_created_at ON chats(created_at);`,

		// Trigger para actualizar updated_at
		`CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql';`,

		`DROP TRIGGER IF EXISTS update_users_updated_at ON users;`,
		`CREATE TRIGGER update_users_updated_at
		BEFORE UPDATE ON users
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();`,

		`DROP TRIGGER IF EXISTS update_chats_updated_at ON chats;`,
		`CREATE TRIGGER update_chats_updated_at
		BEFORE UPDATE ON chats
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();`,

		`DROP TRIGGER IF EXISTS update_admin_groups_updated_at ON admin_groups;`,
		`CREATE TRIGGER update_admin_groups_updated_at
		BEFORE UPDATE ON admin_groups
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();`,

		`DROP TRIGGER IF EXISTS update_workspaces_updated_at ON workspaces;`,
		`CREATE TRIGGER update_workspaces_updated_at
		BEFORE UPDATE ON workspaces
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return err
		}
	}

	log.Println("Migrations completed successfully")
	CreateDefaultAdmin()
	return nil
}

func CreateDefaultAdmin() error {
	// Verificar si ya existe un admin
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Admin user already exists, skipping default admin creation")
		return nil
	}

	// Crear admin por defecto (password: admin123)
	// En producción, esto debería ser cambiado
	hashedPassword, err := utils.HashPassword("manuel10")
	if err != nil {
		log.Printf("Error hashing default admin password: %v", err)
		return err
	}
	_, err = DB.Exec(`
		INSERT INTO users (id, username, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING
	`, uuid.New(), "admin", nil, hashedPassword, "super_admin")

	return err
}
