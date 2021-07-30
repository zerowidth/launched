Rails.application.routes.draw do
  root "plists#new"
  get "help", :controller => "pages", :as => :help
  resources :plists do
    member do
      get :install
    end
  end
end
