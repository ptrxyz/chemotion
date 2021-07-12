# frozen_string_literal: true

require 'rails_helper'

describe 'Collection management' do
  let!(:user) { create(:person) }
  let!(:collection) { create(:collection, user_id: user.id) }

  before do
    user.update!(confirmed_at: Time.now, account_active: true)
    sign_in(user)
  end

  describe 'Createa an unshared collection' do
    it 'create an unshared collection', js: true do
      # press Collections button (on the left-side tree view)
      find('div.take-ownership-btn').click

      # press Add(plus) button to add collection
      find('div.root-actions').find(:xpath, '..').all('button')[0].click

      # input collection name
      factory_collection_name = 'Hello Collection'
      new_collection = find("div#collection-management-tab-pane-0 input[value='New Collection']:last-of-type")
      new_collection.click
      new_collection.set(factory_collection_name)

      # press Update button to save
      find('div.root-actions').find(:xpath, '..').all('button')[0].click

      # except
      expect(find('.tree-view', text: factory_collection_name).text).to eq(factory_collection_name)
    end
  end

  describe 'Delete an unshared collection' do
    it 'delete an unshared collection', js: true do
      # byebug
      # press Collections button (on the left-side tree view)
      find('div.take-ownership-btn').click

      factory_collection_name = collection.label

      # except before deletion
      expect(page).to have_content(factory_collection_name)

      # press Delete button to delete the collection
      find("div#collection-management-tab-pane-0 button[class='btn btn-xs btn-danger']:last-of-type").click

      # press Update button to save
      find('div.root-actions').find(:xpath, '..').all('button')[0].click

      # except after deletion
      expect(page).not_to have_content(factory_collection_name)
    end
  end
end
