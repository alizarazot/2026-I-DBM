use("Plataforma-LINEA");

// 12345678
const passwordHash =
	"bbc47cc6a7c984b9ef6f2f5a9d524346:d380fe5927b32a185c25a66dcc73f74ef9849c0e6a66b2ea8b7e671b16bb5d16832d12286ef851fe96979db35e70e695bbb2e977c4eb1d7dae95042fe5316d3e";

const newUsers = [
	{
		fullName: "Alice Johnson",
		email: "alice.j@example.com",
		role: "student",
		lastName: "Johnson",
		document: "987654321",
		phone: "555-0101",
	},
	{
		fullName: "Robert Smith",
		email: "r.smith@academy.edu",
		role: "teacher",
		lastName: "Smith",
		document: "88223344",
		phone: "555-0102",
	},
	{
		fullName: "Elena Rodriguez",
		email: "elena.admin@local.app",
		role: "manager",
		lastName: "Rodriguez",
		document: "11223344",
		phone: "555-0103",
	},
	{
		fullName: "Kevin Zhang",
		email: "kzhang@student.net",
		role: "student",
		lastName: "Zhang",
		document: "77665544",
		phone: "555-0104",
	},
	{
		fullName: "Sarah Miller",
		email: "s.miller@teacher.com",
		role: "teacher",
		lastName: "Miller",
		document: "55443322",
		phone: "555-0105",
	},
	{
		fullName: "Lucia Fernandez",
		email: "lucia.f@example.org",
		role: "student",
		lastName: "Fernandez",
		document: "99887766",
		phone: "555-0107",
	},
	{
		fullName: "Marcus Brown",
		email: "m.brown@manager.net",
		role: "manager",
		lastName: "Brown",
		document: "55667788",
		phone: "555-0108",
	},
	{
		fullName: "James Taylor",
		email: "jtaylor@academy.edu",
		role: "teacher",
		lastName: "Taylor",
		document: "33445566",
		phone: "555-0110",
	},
	{
		fullName: "Thomas Anderson",
		email: "t.anderson@student.io",
		role: "student",
		lastName: "Anderson",
		document: "10101010",
		phone: "555-0112",
	},
	{
		fullName: "Sophia White",
		email: "swhite@teacher.com",
		role: "teacher",
		lastName: "White",
		document: "12121212",
		phone: "555-0113",
	},
	{
		fullName: "Daniel Lee",
		email: "d.lee@local.app",
		role: "manager",
		lastName: "Lee",
		document: "89898989",
		phone: "555-0114",
	},
	{
		fullName: "Emma Watson",
		email: "ewatson@student.com",
		role: "student",
		lastName: "Watson",
		document: "70707070",
		phone: "555-0115",
	},
	{
		fullName: "Carlos Ruiz",
		email: "cruiz@edu.co",
		role: "teacher",
		lastName: "Ruiz",
		document: "66554433",
		phone: "555-0116",
	},
	{
		fullName: "Hannah Abbott",
		email: "habbott@study.net",
		role: "student",
		lastName: "Abbott",
		document: "55544433",
		phone: "555-0117",
	},
	{
		fullName: "Victor Vance",
		email: "vvance@manager.org",
		role: "manager",
		lastName: "Vance",
		document: "44433322",
		phone: "555-0118",
	},
].map((u) => ({
	...u,
	emailVerified: true,
	banned: false,
	createdAt: new Date(),
	updatedAt: new Date(),
}));

// 1. Insert Users
const userResult = db.getCollection("users").insertMany(newUsers);
const userIds = Object.values(userResult.insertedIds);

// 2. Map Accounts to those User IDs
const accounts = userIds.map((id) => ({
	accountId: id.toString(),
	providerId: "credential",
	userId: id,
	password: passwordHash,
	createdAt: new Date(),
	updatedAt: new Date(),
}));

// 3. Insert Accounts
db.getCollection("user_accounts").insertMany(accounts);

print(
	"Successfully created 15 users and 15 corresponding account credentials.",
);
