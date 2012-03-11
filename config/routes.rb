Launched::Application.routes.draw do
  get "help", :controller => "pages", :as => :help
  resources :plists do
    member do
      get :install
    end
  end
  root :to => "plists#new"
end
