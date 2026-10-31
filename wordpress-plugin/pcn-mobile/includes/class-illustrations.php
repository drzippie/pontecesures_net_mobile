<?php
/**
 * GET /wp-json/pcn/v1/illustrations
 *
 * Publishes the "section → illustration" mapping used by the theme when a
 * post has no image in its body, so the mobile app does not need a hard-coded
 * copy. When the pontecesures theme is active its own functions
 * (inc/imaxes.php) are the source of truth; otherwise the defaults below are
 * returned.
 *
 * @package PcnMobile
 */

namespace PcnMobile;

defined( 'ABSPATH' ) || exit;

final class Illustrations {

	/** Fallback illustration slug. */
	private const FALLBACK = 'xeral';

	public static function register_routes(): void {
		register_rest_route(
			PCN_MOBILE_NAMESPACE,
			'/illustrations',
			array(
				'methods'             => \WP_REST_Server::READABLE,
				'callback'            => array( self::class, 'get' ),
				'permission_callback' => '__return_true',
			)
		);
	}

	public static function get(): \WP_REST_Response {
		$illustrated = function_exists( 'pontecesures_slugs_con_ilustracion' )
			? pontecesures_slugs_con_ilustracion()
			: array( 'sociedade', 'politica', 'cultura', 'festas', 'obras', 'deportes', 'sucesos', 'emprego', 'medio-ambiente', 'esquelas', 'turismo', 'avisos' );

		$equivalents = function_exists( 'pontecesures_equivalencias_categoria' )
			? pontecesures_equivalencias_categoria()
			: array(
				'concello'         => 'politica',
				'pontecesures-cf'  => 'deportes',
				'formacionemprego' => 'emprego',
			);

		$files = function_exists( 'pontecesures_ficheiros_ilustracion' )
			? pontecesures_ficheiros_ilustracion()
			: array(
				'politica'  => 'politica-v2',
				'sociedade' => 'sociedade-v2',
				'turismo'   => 'turismo-v2',
			);

		$base = trailingslashit( get_template_directory_uri() ) . 'assets/img/categorias/';

		$urls = array();
		foreach ( array_merge( $illustrated, array( self::FALLBACK ) ) as $slug ) {
			$name          = $files[ $slug ] ?? $slug;
			$urls[ $slug ] = $base . $name . '.jpg';
		}

		// Equivalents pointing at the fallback are implicit.
		$equivalents = array_filter(
			$equivalents,
			static fn( $target ) => self::FALLBACK !== $target
		);

		$response = new \WP_REST_Response(
			array(
				'version'     => PCN_MOBILE_VERSION,
				'fallback'    => self::FALLBACK,
				'urls'        => $urls,
				'equivalents' => $equivalents,
			)
		);
		// Changes very rarely: let clients and proxies cache it for a day.
		$response->header( 'Cache-Control', 'public, max-age=86400' );
		return $response;
	}
}
