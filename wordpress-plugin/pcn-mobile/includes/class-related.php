<?php
/**
 * GET /wp-json/pcn/v1/posts/{id}/related?per_page=4
 *
 * Returns the IDs of the posts the website shows under "Máis sobre isto",
 * computed by pontecesures_novas_relacionadas() (plugin "comments",
 * includes/relacionadas.php): tags weighted by rarity, a recency bonus and a
 * same-category top-up, cached for 30 minutes per post.
 *
 * Response: {"ids": [123, 456, ...]} in display order. The app then loads the
 * posts through /wp/v2/posts?include=...&orderby=include.
 *
 * If that function is not available the route answers 404 so the app falls
 * back to its own approximation.
 *
 * @package PcnMobile
 */

namespace PcnMobile;

defined( 'ABSPATH' ) || exit;

final class Related {

	private const MAX = 8;

	public static function register_routes(): void {
		register_rest_route(
			PCN_MOBILE_NAMESPACE,
			'/posts/(?P<id>\d+)/related',
			array(
				'methods'             => \WP_REST_Server::READABLE,
				'callback'            => array( self::class, 'get' ),
				'permission_callback' => '__return_true',
				'args'                => array(
					'id'       => array(
						'type'     => 'integer',
						'required' => true,
						'minimum'  => 1,
					),
					'per_page' => array(
						'type'    => 'integer',
						'default' => 4,
						'minimum' => 1,
						'maximum' => self::MAX,
					),
				),
			)
		);
	}

	/**
	 * @param \WP_REST_Request $request Request.
	 * @return \WP_REST_Response|\WP_Error
	 */
	public static function get( \WP_REST_Request $request ) {
		if ( ! function_exists( 'pontecesures_novas_relacionadas' ) ) {
			return new \WP_Error( 'pcn_related_unavailable', 'Related posts are not available.', array( 'status' => 404 ) );
		}

		$id   = (int) $request['id'];
		$post = get_post( $id );
		if ( ! $post || 'post' !== $post->post_type || 'publish' !== $post->post_status ) {
			return new \WP_Error( 'rest_post_invalid_id', 'Invalid post ID.', array( 'status' => 404 ) );
		}

		$posts = pontecesures_novas_relacionadas( $id, (int) $request['per_page'] );
		$ids   = array_values(
			array_map(
				static fn( $p ) => (int) $p->ID,
				array_filter( $posts, static fn( $p ) => $p instanceof \WP_Post && 'publish' === $p->post_status )
			)
		);

		$response = new \WP_REST_Response( array( 'ids' => $ids ) );
		// Same lifetime as the server-side cache of the ranking.
		$response->header( 'Cache-Control', 'public, max-age=1800' );
		return $response;
	}
}
