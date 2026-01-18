DELETE FROM degree_types
WHERE code IN (
               'high_school',
               'vocational',
               'junior_college',
               'bachelor',
               'master',
               'doctor',
               'other'
    );
