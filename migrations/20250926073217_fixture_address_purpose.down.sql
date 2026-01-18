DELETE FROM address_purposes
WHERE code IN (
               'home',
               'contact',
               'office',
               'shipping',
               'billing',
               'other'
              );
