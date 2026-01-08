DELETE FROM education_statuses
WHERE code IN (
               'enrolled',
               'graduated',
               'completed',
               'withdrawn',
               'leave_of_absence'
    );
