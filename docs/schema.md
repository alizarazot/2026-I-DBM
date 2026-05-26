# Schema

The database is MongoDB.

```typescript
// Compound index: `{ role: 1, "basicInfo.genre": 1 }`.
type User = {
    email: string; // Single field index.

    passwordSalt: string;
    passwordHash: string;

    role: "manager" | "teacher" | "student"; // Single field index.

    basicInfo: {
        pictureFileId: ObjectId;

        firstName: string;
        middleName: string;
        firstSurname: string;
        secondSurname: string;

        birthdate: Date;

        genre: "man" | "women" | "other";
    }

    createdAt: Date;
    updatedAt: Date;
};

// TODO: Index can expire JWT token.
type UserSessions = {
    userId: ObjectId;
    jwtToken: string; // Single field index.

    createdAt: Date;
    updatedAt: Date;
};

// Compound index: `{ userId: 1, time: -1 }`.
type UserAnalytics = {
    userId: ObjectId;
    time: Date;

    ip: string;
    location: { // Geospatial Index.
        type: "Point";
        coordinates: [number, number];
    } | null;

    createdAt: Date;
    updatedAt: Date;
};

// CompoundIndex: `{ teacherIds: 1, "schedule.day": 1 }`.
// CompoundIndex: `{ studentIds: 1, "schedule.day": 1 }`.
type Class = {
    name: string;
    desc: string; // Text index.

    maxStudents: number;

    // References the [User] collection.
    teacherIds: []ObjectId;
    studentIds: []OjbectId;

    schedule: {
        day: "sunday" | "monday" | "tuesday" | "wednesday" | "thursday" | "friday" | "saturday";
        startHour: number;
        startMinute: number;
        durationMin: number;
    }[]; // Multikey Index.

    createdAt: Date;
    updatedAt: Date;
};

type Lecture = {
    classId: ObjectId;

    name: string;
    desc: string; // Text index.

    audioFileId: ObjectId;

    startTime: Date;
    endTime: Date;

    createdAt: Date;
    updatedAt: Date;
};

type LectureAssistance = {
    lectureId: ObjectId;
    userAnalyticsIds: ObjectId[];

    createdAt: Date;
    updatedAt: Date;
};

type LectureTranscriptionLine = {
    lectureId: ObjectId;
    userId: ObjectId;

    referenceTranscriptionLineId: ObjectId | null; // In case that someone is answering other's question.
    text: string; // Text index.
    minute: number;

    createdAt: Date;
    updatedAt: Date;
};

type LectureQuestion = {
    lectureId: ObjectId;
    minute: number;

    question: string; // Text index.
    answer: string;

    incorrectAnswers: {
        answer: string;
        teacherSuggestion: string;
    }[];

    results: {
        userId: ObjectId;
        answer: number; // Number -1 means that was correct, else, was one of [incorrectAnswers].
        minute: number;
    }[];

    createdAt: Date;
    updatedAt: Date;
};

// Compound index: `{ userId: 1, createdAt: 1 }`.
// Compound index: `{ category: 1, createdAt: 1 }`.
type CustomerFeedbackAndComplaints = {
    subject: string;
    category: "request" | "complaint" | "claims" | "suggestion";
    userId: ObjectId;

    details: string; // Text index.

    createdAt: Date;
    updatedAt: Date;
};

type CustomerFeedbackAndComplaintsResponse = {
    customerFeedbackAndComplaintId: ObjectId;
    userId: ObjectId; // The user who answers it.

    answer: string; // Text index.

    createdAt: Date;
    updatedAt: Date;
}
```
