import {
  MedicRole,
  TechnicalRole,
  AdministratorRole,
  type EmployeeRole,
  type EmployeeWithDisplayName,
  createDisplayName,
  isAdmin,
  isMedic,
  isTechnical
} from './employee-extensions';
import { EmployeeResponse } from '../index';
import { EmployeeResponseProfileTypeEnum } from '../index';

describe('Employee Extensions', () => {
  let mockEmployee: EmployeeResponse;

  beforeEach(() => {
    mockEmployee = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: EmployeeResponseProfileTypeEnum.Medic,
      username: 'johndoe',
      gender: 'Male'
    };
  });

  describe('Role Constants', () => {
    it('should have correct role constants', () => {
      expect(MedicRole).toBe('Medic');
      expect(TechnicalRole).toBe('Technical');
      expect(AdministratorRole).toBe('Administrator');
    });
  });

  describe('createDisplayName', () => {
    it('should create display name from first and last name', () => {
      const displayName = createDisplayName(mockEmployee);
      expect(displayName).toBe('John Doe');
    });

    it('should handle empty last name', () => {
      mockEmployee.lastName = '';
      const displayName = createDisplayName(mockEmployee);
      expect(displayName).toBe('John');
    });

    it('should handle empty first name', () => {
      mockEmployee.firstName = '';
      const displayName = createDisplayName(mockEmployee);
      expect(displayName).toBe('Doe');
    });

    it('should handle both names empty', () => {
      mockEmployee.firstName = '';
      mockEmployee.lastName = '';
      const displayName = createDisplayName(mockEmployee);
      expect(displayName).toBe('');
    });

    it('should trim whitespace', () => {
      mockEmployee.firstName = '  John  ';
      mockEmployee.lastName = '  Doe  ';
      const displayName = createDisplayName(mockEmployee);
      expect(displayName).toBe('John     Doe');
    });
  });

  describe('isAdmin', () => {
    it('should return true for Administrator role', () => {
      mockEmployee.profileType = EmployeeResponseProfileTypeEnum.Administrator;
      expect(isAdmin(mockEmployee)).toBe(true);
    });

    it('should return false for Medic role', () => {
      mockEmployee.profileType = EmployeeResponseProfileTypeEnum.Medic;
      expect(isAdmin(mockEmployee)).toBe(false);
    });

    it('should return false for Technical role', () => {
      mockEmployee.profileType = EmployeeResponseProfileTypeEnum.Technical;
      expect(isAdmin(mockEmployee)).toBe(false);
    });

    it('should return false for unknown role', () => {
      mockEmployee.profileType = 'Unknown' as any;
      expect(isAdmin(mockEmployee)).toBe(false);
    });
  });

  describe('isMedic', () => {
    it('should return true for Medic role', () => {
      mockEmployee.profileType = EmployeeResponseProfileTypeEnum.Medic;
      expect(isMedic(mockEmployee)).toBe(true);
    });

    it('should return false for Administrator role', () => {
      mockEmployee.profileType = EmployeeResponseProfileTypeEnum.Administrator;
      expect(isMedic(mockEmployee)).toBe(false);
    });

    it('should return false for Technical role', () => {
      mockEmployee.profileType = EmployeeResponseProfileTypeEnum.Technical;
      expect(isMedic(mockEmployee)).toBe(false);
    });

    it('should return false for unknown role', () => {
      mockEmployee.profileType = 'Unknown' as any;
      expect(isMedic(mockEmployee)).toBe(false);
    });
  });

  describe('isTechnical', () => {
    it('should return true for Technical role', () => {
      mockEmployee.profileType = EmployeeResponseProfileTypeEnum.Technical;
      expect(isTechnical(mockEmployee)).toBe(true);
    });

    it('should return false for Administrator role', () => {
      mockEmployee.profileType = EmployeeResponseProfileTypeEnum.Administrator;
      expect(isTechnical(mockEmployee)).toBe(false);
    });

    it('should return false for Medic role', () => {
      mockEmployee.profileType = EmployeeResponseProfileTypeEnum.Medic;
      expect(isTechnical(mockEmployee)).toBe(false);
    });

    it('should return false for unknown role', () => {
      mockEmployee.profileType = 'Unknown' as any;
      expect(isTechnical(mockEmployee)).toBe(false);
    });
  });

  describe('EmployeeWithDisplayName interface', () => {
    it('should extend EmployeeResponse with displayName', () => {
      const employeeWithDisplayName: EmployeeWithDisplayName = {
        ...mockEmployee,
        displayName: createDisplayName(mockEmployee)
      };

      expect(employeeWithDisplayName.displayName).toBe('John Doe');
      expect(employeeWithDisplayName.firstName).toBe('John');
      expect(employeeWithDisplayName.lastName).toBe('Doe');
      expect(employeeWithDisplayName.email).toBe('john.doe@example.com');
    });
  });

  describe('EmployeeRole type', () => {
    it('should accept valid role values', () => {
      const medicRole: EmployeeRole = MedicRole;
      const technicalRole: EmployeeRole = TechnicalRole;
      const adminRole: EmployeeRole = AdministratorRole;

      expect(medicRole).toBe('Medic');
      expect(technicalRole).toBe('Technical');
      expect(adminRole).toBe('Administrator');
    });
  });
});
