// DO NOT EDIT. This file is automatically generated.

export interface Status {
	email: EMail;
	roles: Roles;
}

export interface LoginStatus {
	Email: string;
	ID: string;
	LoginURL: string;
	IsAGoogler: boolean;
	IsAdmin: boolean;
	IsEditor: boolean;
	IsViewer: boolean;
}

export type EMail = string;

export type Role = 'viewer' | 'editor' | 'admin' | '';

export type Roles = Role[] | null;
