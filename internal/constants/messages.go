package constants

const (
	// ==================================
	// Éxito Archivos
	// ==================================
	MsgImageUploadedSuccessfully = "imagen subida exitosamente"
	MsgVideoUploadedSuccessfully = "video subido exitosamente"

	// ==================================
	// Éxito Base de Datos
	// ==================================
	MsgDBConnectionSuccess = "conexión a la base de datos exitosa"
	MsgDBMigrationSuccess  = "migraciones de base de datos ejecutadas correctamente"
	MsgDBSuccess           = "base de datos inicializada correctamente"

	// ==================================
	// Éxito Cloudinary
	// ==================================
	MsgCloudinarySuccess = "cloudinary inicializado correctamente"

	// ==================================
	// Éxito Películas
	// ==================================
	MsgMovieCreatedSuccessfully    = "película creada exitosamente"
	MsgMovieDeletedSuccessfully    = "película eliminada exitosamente"
	MsgMovieRetrievedSuccessfully  = "película obtenida exitosamente"
	MsgMovieUpdatedSuccessfully    = "película actualizada exitosamente"
	MsgMoviesRetrievedSuccessfully = "películas obtenidas correctamente"
	MsgMovieDeletedLog             = "solicitud de eliminación de película con ID"
	MsgMovieUpdatedFieldsLog       = "campos específicos de película actualizados para ID"
	MsgMovieNotFound               = "película no encontrada con ID"

	// ==================================
	// Éxito Perfiles
	// ==================================
	MsgProfileCreatedSuccessfully    = "Perfil creado con éxito para el usuario ID"
	MsgProfileDeletedSuccessfully    = "Perfil eliminado con éxito para el usuario ID"
	MsgProfileRetrievedSuccessfully  = "Perfil obtenido con éxito para el usuario ID"
	MsgProfileUpdatedSuccessfully    = "Perfil actualizado con éxito para el usuario ID"
	MsgProfilesRetrievedSuccessfully = "Perfiles obtenidos con éxito para el usuario ID"

	// ==================================
	// Éxito Servidor
	// ==================================
	MsgEnvCargedSuccessfully     = "variables de entorno cargadas correctamente"
	MsgServerStartedSuccessfully = "servidor iniciado correctamente"

	// ==================================
	// Éxito Usuarios
	// ==================================
	MsgLoginSuccessful            = "inicio de sesión exitoso"
	MsgUserCreatedSuccessfully    = "usuario creado exitosamente en base de datos"
	MsgUserFoundByEmail           = "usuario encontrado por email"
	MsgUserLoggedInSuccessfully   = "usuario autenticado correctamente"
	MsgUserRegisteredSuccessfully = "usuario registrado correctamente"

	// ==================================
	// Errores de Archivos/Cloudinary
	// ==================================
	ErrMsgInitCloudinary      = "error al inicializar cloudinary"
	ErrMsgOpenImage           = "error al abrir la imagen"
	ErrMsgOpenUploadedFile    = "no se pudo abrir el archivo subido"
	ErrMsgOpenVideo           = "error al abrir el video"
	ErrMsgUnsupportedFileType = "tipo de archivo no soportado"
	ErrMsgUploadImage         = "error al subir la imagen"
	ErrMsgUploadToCloudinary  = "error al subir el archivo a cloudinary"
	ErrMsgUploadVideo         = "error al subir el video"

	// ==================================
	// Errores de Autenticación/Autorización
	// ==================================
	ErrAdminOnly               = "requiere permisos de administrador"
	ErrInvalidSigningMethod    = "algoritmo de firma inválido"
	ErrInvalidToken            = "token inválido"
	ErrMissingOrMalformedToken = "token faltante o malformado"

	// ==================================
	// Errores de Base de Datos
	// ==================================
	ErrMsgCreateMovie         = "error al crear la película"
	ErrMsgCreateUser          = "error al crear el usuario en la base de datos"
	ErrMsgDBAccessDenied      = "acceso denegado a la base de datos"
	ErrMsgDBConnection        = "no se pudo conectar a la base de datos"
	ErrMsgDBMigration         = "error al ejecutar las migraciones"
	ErrMsgDeleteMovie         = "error al eliminar la película"
	ErrMsgDeleteProfile       = "error al eliminar el perfil"
	ErrMsgFindUserByEmail     = "error al buscar el usuario por email"
	ErrMsgGetMovieByID        = "error al obtener la película por ID"
	ErrMsgGetMovies           = "error al obtener las películas"
	ErrMsgGetProfile          = "error al obtener el perfil"
	ErrMsgGetProfiles         = "error al obtener los perfiles"
	ErrMsgGetUpdatedMovie     = "error al obtener la película actualizada"
	ErrMsgProfileNotFound     = "perfil no encontrado o no permitido"
	ErrMsgRegisterProfile     = "no se pudo registrar el perfil"
	ErrMsgRegisterUser        = "no se pudo registrar el usuario"
	ErrMsgUpdateMovie         = "error al actualizar la película"
	ErrMsgUpdateProfile       = "error al actualizar el perfil"
	ErrMsgUserNotFound        = "inicio de sesión fallido - usuario no encontrado"
	ErrMsgUserNotFoundByEmail = "usuario no encontrado con ese email"
	ErrMsgUpdateMovieFieldsDB = "error al actualizar campos de la película"
	ErrMsgGetMovieAll         = "error al obtener todas las películas"
	ErrMsgDeleteMovieByID     = "error al eliminar la película con ID"
	ErrMsgMovieNotFound       = "película no encontrada"
	ErrMsgGetMovie            = "error al obtener la película"

	// ==================================
	// Errores de Validación
	// ==================================
	ErrMsgEmailAlreadyRegistered = "el correo ya está registrado"
	ErrMsgInvalidCredentials     = "credenciales incorrectas"
	ErrMsgInvalidDuration        = "duración no válida"
	ErrMsgInvalidID              = "id no válido"
	ErrMsgInvalidRequest         = "datos inválidos"
	ErrMsgInvalidYear            = "año no válido"
	ErrMsgInvalidLoginData       = "datos de inicio de sesión inválidos"
	ErrMsgUpdateMovieFields      = "campos inválidos para actualizar la película"
	ErrMsgEmptyEmail             = "el email no puede estar vacío"
	ErrMsgEmptyPassword          = "la contraseña no puede estar vacía"
	ErrMsgInvalidRegisterData    = "datos de registro inválidos"
	ErrMshEmptyProfileName       = "el nombre de perfil no puede estar vacío"
	ErrMsgEmptyProfileId         = "el id del perfil no puede estar vacío"
	ErrMsgInvalidProfileName     = "el nombre del perfil debe tener entre 1 y 10 caracteres"
	ErrMsgInvalidAvatarURL       = "la URL del avatar no es válida"
	ErrMsgMaxProfilesReached     = "ya tienes el número máximo de perfiles (4)"
	ErrMsgProfileNameTaken       = "ya existe un perfil con ese nombre"
	ErrMsgInvalidInput           = "datos de entrada inválidos"

	// ==================================
	// Errores Técnicos/Generales
	// ==================================
	ErrMsgGeneratingToken     = "error generando token jwt"
	ErrMsgHashingPassword     = "error al hashear la contraseña"
	ErrMsgInternalServerError = "error interno del servidor"
	ErrMsgMissingEnvVar       = "variable de entorno faltante"
	ErrMsgParseForm           = "error al parsear el formulario"
	ErrMsgServerStart         = "error al iniciar el servidor"
)
