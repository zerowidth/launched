Launched::Application.routes.draw do
  resources :plists do
    member do
      get :install
    end
  end
  root :to => "plists#new"
end
