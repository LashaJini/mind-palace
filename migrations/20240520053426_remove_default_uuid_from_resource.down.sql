ALTER TABLE {{ .Namespace }}.resource ALTER COLUMN id SET DEFAULT gen_random_uuid();
