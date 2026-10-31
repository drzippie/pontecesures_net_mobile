<?php
/**
 * Plugin Name:       Pontecesures.net Mobile API
 * Description:       REST endpoints used by the Pontecesures.net mobile app (namespace pcn/v1).
 * Version:           0.1.0
 * Requires at least: 6.4
 * Requires PHP:      8.1
 * Author:            Antonio Cortés (DrZippie)
 * License:           MIT
 * License URI:       https://opensource.org/licenses/MIT
 * Text Domain:       pcn-mobile
 *
 * @package PcnMobile
 */

defined( 'ABSPATH' ) || exit;

const PCN_MOBILE_VERSION   = '0.1.0';
const PCN_MOBILE_NAMESPACE = 'pcn/v1';

require_once __DIR__ . '/includes/class-illustrations.php';

add_action(
	'rest_api_init',
	static function () {
		PcnMobile\Illustrations::register_routes();
		// Authentication and registration endpoints (pcn/v1/auth/*) will be
		// added here in a later phase. See README.md.
	}
);
