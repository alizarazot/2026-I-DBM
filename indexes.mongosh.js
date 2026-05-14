use('Plataforma-LINEA');

// 1. Single Field Indexes
db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ role: 1 });
db.user_sessions.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0 });
db.user_sessions.createIndex({ token: 1 }, { unique: true });
db.transcriptions.createIndex({ sessionId: 1 });

// 2. Compound Index
db.users.createIndex({ lastName: 1, firstName: 1 });

// 3. Multikey Index (indexing the studentsIds array in courses)
db.courses.createIndex({ studentsIds: 1 });

// 4. Text index.
db.courses.createIndex({ name: 'text_course', name: 'text', description: 'text' });

// 5. Geospatial Index
db.user_location_analytics.createIndex({ coordinates: '2dsphere' });
