import { InMemoryDbService } from 'angular-in-memory-web-api';

export class EmployeeData implements InMemoryDbService {
  createDb() {
    const employees = [
      { id: 1, username: "jdoe", firstName: "John", lastName: "Doe", role: "Medic", email: "jdoe@example.com" },
      { id: 2, username: "asmith", firstName: "Alice", lastName: "Smith", role: "Technical", email: "asmith@example.com" },
    ];
    return { employees };
  }
}
