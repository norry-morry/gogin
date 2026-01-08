DELETE FROM genders
WHERE code IN (
               'male',
               'female',
               'other',
               'unspecified'
    );
