use('Plataforma-LINEA');

db.users.createIndex({ email: 1 }, { unique: true });

// Role is kinda 'sorter index', this case ascending.
db.users.createIndex({ role: 1 });

//db.courses.createIndex({ teacherId: 1 });

// Token needs to be validated on each request.
db.user_sessions.createIndex({ expiresAt: 1 });
db.user_sessions.createIndex({ token: 1 }, { unique: true });

db.transcriptions.createIndex({ sessionId: 1 });
