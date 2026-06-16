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
