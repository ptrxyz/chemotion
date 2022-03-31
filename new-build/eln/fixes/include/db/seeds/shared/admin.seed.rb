attributes = {
  email: 'eln-admin@kit.edu',
  first_name: 'ELN',
  last_name: 'Admin',
  password: 'PleaseChangeYourPassword',
  name_abbreviation: 'ADM',
  type: 'Admin',
  account_active: true
}

User.create!(attributes) unless User.find_by(type: attributes[:type], name_abbreviation: attributes[:name_abbreviation])

attributes = {
  email: 'eln-user@kit.edu',
  first_name: 'ELN',
  last_name: 'User',
  password: 'chemotion',
  name_abbreviation: 'CU1',
  type: 'Person',
  account_active: true 
}

User.create!(attributes) unless User.find_by(type: attributes[:type], name_abbreviation: attributes[:name_abbreviation])
