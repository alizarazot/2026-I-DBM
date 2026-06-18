export interface User {
	email: string;
	role: UserRole;
	info: UserInfo;
}

export type UserRole = 'invalid' | 'manager' | 'teacher' | 'student';

export interface UserInfo {
	firstName: string;
	middleName: string;
	firstSurname: string;
	secondSurname: string;
	birthdate: Date;
	genre: UserGenre;
}

export type UserGenre = 'invalid' | 'male' | 'female' | 'other';

export type CfcCategory = 'request' | 'complaint' | 'claim' | 'suggestion' | 'invalid';

export interface Cfc {
	id: string;
	subject: string;
	category: CfcCategory;
	userEmail: string;
	details: string;
	answered: boolean;
	updatedAt: Date;
}

export interface CfcAnswer {
	id: string;
	cfcId: string;
	userEmail: string;
	answer: string;
	updatedAt: string;
}
