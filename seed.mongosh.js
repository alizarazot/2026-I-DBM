use('Plataforma-LINEA');

// 12345678
const passwordHash =
	'bbc47cc6a7c984b9ef6f2f5a9d524346:d380fe5927b32a185c25a66dcc73f74ef9849c0e6a66b2ea8b7e671b16bb5d16832d12286ef851fe96979db35e70e695bbb2e977c4eb1d7dae95042fe5316d3e';

// --- PART 1: USERS ---
const newUsers = [
	{
		firstName: 'Alice',
		email: 'alice.j@example.com',
		role: 'student',
		lastName: 'Johnson',
		document: '987654321',
		phone: '555-0101'
	},
	{
		firstName: 'Robert',
		email: 'r.smith@academy.edu',
		role: 'teacher',
		lastName: 'Smith',
		document: '88223344',
		phone: '555-0102'
	},
	{
		firstName: 'Elena',
		email: 'elena.admin@local.app',
		role: 'manager',
		lastName: 'Rodriguez',
		document: '11223344',
		phone: '555-0103'
	},
	{
		firstName: 'Kevin',
		email: 'kzhang@student.net',
		role: 'student',
		lastName: 'Zhang',
		document: '77665544',
		phone: '555-0104'
	},
	{
		firstName: 'Sarah',
		email: 's.miller@teacher.com',
		role: 'teacher',
		lastName: 'Miller',
		document: '55443322',
		phone: '555-0105'
	},
	{
		firstName: 'Lucia',
		email: 'lucia.f@example.org',
		role: 'student',
		lastName: 'Fernandez',
		document: '99887766',
		phone: '555-0107'
	},
	{
		firstName: 'Marcus',
		email: 'm.brown@manager.net',
		role: 'manager',
		lastName: 'Brown',
		document: '55667788',
		phone: '555-0108'
	},
	{
		firstName: 'James',
		email: 'jtaylor@academy.edu',
		role: 'teacher',
		lastName: 'Taylor',
		document: '33445566',
		phone: '555-0110'
	},
	{
		firstName: 'Thomas',
		email: 't.anderson@student.io',
		role: 'student',
		lastName: 'Anderson',
		document: '10101010',
		phone: '555-0112'
	},
	{
		firstName: 'Sophia',
		email: 'swhite@teacher.com',
		role: 'teacher',
		lastName: 'White',
		document: '12121212',
		phone: '555-0113'
	},
	{
		firstName: 'Daniel',
		email: 'd.lee@local.app',
		role: 'manager',
		lastName: 'Lee',
		document: '89898989',
		phone: '555-0114'
	},
	{
		firstName: 'Emma',
		email: 'ewatson@student.com',
		role: 'student',
		lastName: 'Watson',
		document: '70707070',
		phone: '555-0115'
	},
	{
		firstName: 'Carlos',
		email: 'cruiz@edu.co',
		role: 'teacher',
		lastName: 'Ruiz',
		document: '66554433',
		phone: '555-0116'
	},
	{
		firstName: 'Hannah',
		email: 'habbott@study.net',
		role: 'student',
		lastName: 'Abbott',
		document: '55544433',
		phone: '555-0117'
	},
	{
		firstName: 'Victor',
		email: 'vvance@manager.org',
		role: 'manager',
		lastName: 'Vance',
		document: '44433322',
		phone: '555-0118'
	}
].map((u) => ({
	...u,
	emailVerified: true,
	banned: false,
	createdAt: new Date(),
	updatedAt: new Date()
}));

const userResult = db.getCollection('users').insertMany(newUsers);
const userIds = Object.values(userResult.insertedIds);

// --- PART 2: ACCOUNTS ---
const accounts = userIds.map((id) => ({
	accountId: id.toString(),
	providerId: 'credential',
	userId: id,
	password: passwordHash,
	createdAt: new Date(),
	updatedAt: new Date()
}));
db.getCollection('user_accounts').insertMany(accounts);

// --- PART 3: COURSES ---
// Get the IDs of the teachers we just created to assign them to courses
const teachers = db.getCollection('users').find({ role: 'teacher' }).toArray();
const teacherIds = teachers.map((t) => t._id.toString());

const daysOfWeek = ['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday'];

const courseData = [
	{
		name: 'Introducción a Programación',
		description: 'Conceptos básicos de lógica y algoritmos.',
		maxStudents: 30
	},
	{
		name: 'Matemáticas Avanzadas',
		description: 'Cálculo diferencial e integral.',
		maxStudents: 25
	},
	{
		name: 'Historia Universal',
		description: 'Un recorrido por los hitos de la humanidad.',
		maxStudents: 40
	},
	{
		name: 'Diseño Gráfico Digital',
		description: 'Uso de herramientas para creación visual.',
		maxStudents: 20
	},
	{
		name: 'Inglés Técnico',
		description: 'Vocabulario especializado para profesionales.',
		maxStudents: 35
	},
	{
		name: 'Gestión de Proyectos',
		description: 'Metodologías ágiles y tradicionales.',
		maxStudents: 15
	},
	{
		name: 'Base de Datos NoSQL',
		description: 'Modelado y consulta en MongoDB.',
		maxStudents: 25
	},
	{
		name: 'Marketing Digital',
		description: 'Estrategias de SEO y redes sociales.',
		maxStudents: 50
	},
	{
		name: 'Física Cuántica',
		description: 'Introducción a la mecánica cuántica.',
		maxStudents: 10
	},
	{
		name: 'Escritura Creativa',
		description: 'Taller de narrativa y poesía.',
		maxStudents: 15
	},
	{
		name: 'Desarrollo Web Svelte',
		description: 'Creación de apps modernas con Svelte 5.',
		maxStudents: 20
	},
	{
		name: 'Inteligencia Artificial',
		description: 'Fundamentos de redes neuronales.',
		maxStudents: 12
	},
	{
		name: 'Sistemas Operativos',
		description: 'Arquitectura y administración de Linux.',
		maxStudents: 30
	},
	{
		name: 'Economía Aplicada',
		description: 'Análisis de mercados y finanzas.',
		maxStudents: 45
	},
	{
		name: 'Psicología Organizacional',
		description: 'Liderazgo y comportamiento humano.',
		maxStudents: 25
	}
].map((course, index) => {
	// Generar hora militar aleatoria (HH:mm)
	const randomHour = String(Math.floor(Math.random() * 24)).padStart(2, '0');
	const randomMinute = String(Math.floor(Math.random() * 60)).padStart(2, '0');

	return {
		...course,
		teacherId: teacherIds[index % teacherIds.length],
		// Asigna día rotativo (excluyendo domingo)
		day: daysOfWeek[index % daysOfWeek.length],
		// Formato HH:mm
		startHour: `${randomHour}:${randomMinute}`,
		// Entero aleatorio entre 1 y 5
		duration: Math.floor(Math.random() * 5) + 1
	};
});

db.getCollection('courses').insertMany(courseData);
