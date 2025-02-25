use actix_web::{web, App, HttpResponse, HttpServer, Responder};
use actix_web::dev::{ServiceRequest};
use actix_web::middleware::Logger;
use chrono::Utc;
use jsonwebtoken::{encode, decode, Header, Validation, EncodingKey, DecodingKey, Algorithm, TokenData};
use serde::{Serialize, Deserialize};
use std::collections::HashSet;
use std::sync::Mutex;
use actix_web::dev::Service; // Make sure the Service trait is in scope

#[derive(Debug, Serialize, Deserialize)]
struct Claims {
    api_key: String,
    exp: usize,
}

#[derive(Debug, Serialize, Deserialize)]
struct AuthRequest {
    apiKey: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct AuthResponse {
    token: String,
}

struct AppState {
    valid_api_keys: Mutex<HashSet<String>>,
    secret: String,
}

// Handler for issuing tokens
async fn auth_token(
    data: web::Data<AppState>,
    req: web::Json<AuthRequest>,
) -> impl Responder {
    let keys = data.valid_api_keys.lock().unwrap();
    if !keys.contains(&req.apiKey) {
        return HttpResponse::Unauthorized().json("Invalid API key");
    }
    drop(keys);

    let expiration = Utc::now()
        .checked_add_signed(chrono::Duration::hours(1))
        .expect("valid timestamp")
        .timestamp() as usize;
    let claims = Claims {
        api_key: req.apiKey.clone(),
        exp: expiration,
    };

    let token = encode(
        &Header::new(Algorithm::HS256),
        &claims,
        &EncodingKey::from_secret(data.secret.as_ref()),
    ).unwrap();

    HttpResponse::Ok().json(AuthResponse { token })
}

// Validate JWT token from a token string
fn validate_token(data: &web::Data<AppState>, token: &str) -> Result<TokenData<Claims>, jsonwebtoken::errors::Error> {
    decode::<Claims>(
        token,
        &DecodingKey::from_secret(data.secret.as_ref()),
        &Validation::new(Algorithm::HS256),
    )
}

// Protected endpoint handler
async fn protected_endpoint() -> impl Responder {
    HttpResponse::Ok().json("This is protected data accessible only with a valid token.")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Initialize valid API keys and secret
    let mut keys = HashSet::new();
    keys.insert("test_key".to_string());
    let app_state = web::Data::new(AppState {
        valid_api_keys: Mutex::new(keys),
        secret: "your-secret-key".to_string(),
    });

    HttpServer::new(move || {
        let state = app_state.clone();
        App::new()
            .app_data(state.clone())
            .wrap(Logger::default())
            // Public endpoint for token issuance
            .service(
                web::resource("/auth/token")
                    .route(web::post().to(auth_token))
            )
            // Protected endpoints with authentication middleware using inline wrap_fn
            .service(
                web::scope("")
                    .wrap_fn(move |req: ServiceRequest, mut srv| {
                        let state_inner = state.clone();
                        async move {
                            if let Some(auth_header) = req.headers().get("authorization") {
                                if let Ok(auth_str) = auth_header.to_str() {
                                    if auth_str.starts_with("Bearer ") {
                                        let token = &auth_str[7..];
                                        if validate_token(&state_inner, token).is_ok() {
                                            return srv.call(req).await;
                                        } else {
                                            let response = HttpResponse::Unauthorized()
                                                .json("Invalid or expired token");
                                            return Ok(req.into_response(response.map_into_boxed_body()));
                                        }
                                    }
                                }
                            }
                            let response = HttpResponse::Unauthorized()
                                .json("Missing or invalid Authorization header");
                            Ok(req.into_response(response.map_into_boxed_body()))
                        }
                    })
                    .service(
                        web::resource("/protected")
                            .route(web::get().to(protected_endpoint))
                    )
            )
    })
    .bind("127.0.0.1:3000")?
    .run()
    .await
}